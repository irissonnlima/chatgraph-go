// Package service provides the main chatbot application service.
// This file contains the Engine, which handles route registration and execution logic.
package service

import (
	"context"
	"fmt"
	"log"
	"regexp"

	route_return "github.com/irissonnlima/chatgraph-go/core/domain"
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_context "github.com/irissonnlima/chatgraph-go/core/domain/context"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_router "github.com/irissonnlima/chatgraph-go/core/domain/router"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
	adapter_output "github.com/irissonnlima/chatgraph-go/core/ports/adapters/output"
)

// Engine handles route registration and execution logic without I/O.
// It is designed to be testable in isolation without requiring external dependencies.
type Engine[Obs any] struct {
	// routes maps route names to their handler administrators.
	routes map[string]d_router.RouterHandlerAdmnistrator[Obs]
	// defaultOptions holds the default handler options.
	defaultOptions d_router.RouterHandlerOptions
	// routeTriggers holds the route triggers configuration.
	routeTriggers []d_router.RouteTrigger
}

// NewEngine creates a new Engine instance with optional default options.
//
// If not provided, the following defaults are used:
//   - Timeout: 5 minutes (redirects to "timeout_route")
//   - Loop Limit: 3 iterations (redirects to "loop_route")
//   - Protected: nil (no protection by default)
func NewEngine[Obs any](defaultOptions ...d_router.RouterHandlerOptions) *Engine[Obs] {
	defaultOpts := d_router.RouterHandlerOptions{
		Timeout:   &d_router.DEFAULT_TIMEOUT,
		LoopCount: &d_router.DEFAULT_LOOP_COUNT,
		Protected: nil,
	}

	if len(defaultOptions) > 0 {
		d := defaultOptions[0]
		defaultOpts.SetOps(d)
	}

	return &Engine[Obs]{
		routes:         make(map[string]d_router.RouterHandlerAdmnistrator[Obs]),
		defaultOptions: defaultOpts,
	}
}

// RegisterRoute registers a route handler with optional configuration.
// The route name is case-sensitive and must be unique.
func (e *Engine[Obs]) RegisterRoute(
	route string,
	handler d_router.RouteHandler[Obs],
	options ...d_router.RouterHandlerOptions,
) {
	rho := d_router.RouterHandlerOptions{
		Timeout:   e.defaultOptions.Timeout,
		LoopCount: e.defaultOptions.LoopCount,
		Protected: e.defaultOptions.Protected,
	}

	if len(options) > 0 {
		opt := options[0]
		rho.SetOps(opt)
	}

	e.routes[route] = d_router.RouterHandlerAdmnistrator[Obs]{
		HandlerOptions: rho,
		Handler:        handler,
	}
}

// RegisterTrigger registers a global trigger that applies to all routes.
// Triggers are regex patterns that, when matched, redirect to a specific route.
func (e *Engine[Obs]) RegisterTrigger(trigger d_router.RouteTrigger) {
	e.routeTriggers = append(e.routeTriggers, trigger)
}

// applyTriggers checks if the message matches any trigger regex.
// If a match is found, returns the route associated with the trigger.
// Returns empty string if no trigger matches.
func (e *Engine[Obs]) applyTriggers(messageText string) string {
	for _, trigger := range e.routeTriggers {
		re, err := regexp.Compile(trigger.Regex)
		if err != nil {
			log.Printf("[ERROR] invalid trigger regex: %s - %v", trigger.Regex, err)
			continue
		}

		if re.MatchString(messageText) {
			return trigger.Route
		}
	}
	return ""
}

// Execute processes a message using the provided router.
// This is the core method for executing route handlers.
//
// It handles:
//   - Trigger matching
//   - Loop detection
//   - Handler execution with timeout
func (e *Engine[Obs]) Execute(
	userState d_user.UserState[Obs],
	message d_message.Message,
	router adapter_output.IBotExecutor,
) (route_return.RouteReturn, error) {
	// Check for triggers
	preRoute := e.applyTriggers(message.EntireText())
	route := userState.Route

	if preRoute != "" && preRoute != route.Current() {
		log.Printf("[INFO] Triggered route change to: %s", preRoute)
		return &d_action.RedirectResponse{TargetRoute: preRoute}, nil
	}

	// Check for loops
	loopLimit := e.defaultOptions.LoopCount.Count
	repeated := route.CurrentRepeated()
	if repeated > loopLimit && route.Current() != e.defaultOptions.LoopCount.Route {
		log.Printf("[ERROR] Loop detected for route: %s", route.Current())
		return &d_action.RedirectResponse{
			TargetRoute: e.defaultOptions.LoopCount.Route,
		}, nil
	}

	// Get route handler
	routeFunc, exists := e.routes[route.Current()]
	if !exists {
		return nil, fmt.Errorf("route not found: %s", route.Current())
	}

	// Create context with router
	ctx, cancel := d_context.NewChatContext(
		userState,
		message,
		router,
		routeFunc.HandlerOptions.Timeout.Duration,
	)
	defer cancel()

	// Channel to receive the result
	resultChan := make(chan route_return.RouteReturn, 1)

	// Execute handler in goroutine
	go func() {
		resultChan <- routeFunc.Handler(&ctx)
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		if result == nil {
			return userState.Route.Next(userState.Route.Current()), nil
		}
		return result, nil

	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("[ERROR] Handler timeout for route: %s", route.Current())
			return &d_action.RedirectResponse{
				TargetRoute: routeFunc.HandlerOptions.Timeout.Route,
			}, nil
		}
		return nil, ctx.Err()
	}
}

// ValidateRoutes checks that all required routes are registered.
// Returns an error if validation fails.
func (e *Engine[Obs]) ValidateRoutes() error {
	// Check if "start" route exists
	if _, exists := e.routes["start"]; !exists {
		return fmt.Errorf("required route 'start' is not registered")
	}

	// Check if all trigger routes exist
	for _, trigger := range e.routeTriggers {
		if _, exists := e.routes[trigger.Route]; !exists {
			return fmt.Errorf("trigger route '%s' (regex: %s) is not registered", trigger.Route, trigger.Regex)
		}
	}

	// Also check triggers defined in individual route options
	for routeName, handler := range e.routes {
		for _, trigger := range handler.HandlerOptions.Triggers {
			if _, exists := e.routes[trigger.Route]; !exists {
				return fmt.Errorf("trigger route '%s' in route '%s' (regex: %s) is not registered",
					trigger.Route, routeName, trigger.Regex)
			}
		}
	}

	// Also check the default options triggers
	rhoRoutes := e.defaultOptions.GetRhoRoutes()
	for _, rhoRoute := range rhoRoutes {
		if _, exists := e.routes[rhoRoute]; !exists {
			return fmt.Errorf("default option route '%s' is not registered", rhoRoute)
		}
	}

	log.Printf("[INFO] Routes validated successfully. %d routes registered.", len(e.routes))
	return nil
}
