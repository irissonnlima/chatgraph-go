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
	d_file "github.com/irissonnlima/chatgraph-go/core/domain/file"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_route "github.com/irissonnlima/chatgraph-go/core/domain/route"
	d_router "github.com/irissonnlima/chatgraph-go/core/domain/router"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
	adapter_output "github.com/irissonnlima/chatgraph-go/core/ports/adapters/output"
)

// ExecuteResult contains the result of executing a route handler.
// It captures all outputs that would normally be sent to external services.
type ExecuteResult[Obs any] struct {
	// Messages contains all messages that were sent during handler execution.
	Messages []d_message.Message
	// NextRoute is the route set for the user's next message (empty if redirect/end).
	NextRoute string
	// Redirect is set if a redirect was triggered (immediate route change).
	Redirect *d_action.RedirectResponse
	// EndSession is set if the session should be ended.
	EndSession *d_action.EndAction
	// Transfer is set if the user should be transferred to another menu.
	Transfer *d_action.TransferToMenu
	// Observation is the updated observation after handler execution.
	Observation Obs
	// ObservationChanged indicates if the observation was modified.
	ObservationChanged bool
	// Error contains any error that occurred during execution.
	Error error
	// TimedOut indicates if the handler timed out.
	TimedOut bool
}

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

// GetRoutes returns a copy of the registered routes map.
func (e *Engine[Obs]) GetRoutes() map[string]d_router.RouterHandlerAdmnistrator[Obs] {
	return e.routes
}

// GetDefaultOptions returns the default handler options.
func (e *Engine[Obs]) GetDefaultOptions() d_router.RouterHandlerOptions {
	return e.defaultOptions
}

// GetTriggers returns the registered triggers.
func (e *Engine[Obs]) GetTriggers() []d_router.RouteTrigger {
	return e.routeTriggers
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

// mockRouter is an internal mock for capturing messages during execution.
type mockRouter[Obs any] struct {
	messages           []d_message.Message
	observation        string
	observationChanged bool
	nextRoute          string
	endSession         *d_action.EndAction
	transfer           *d_action.TransferToMenu
}

func (m *mockRouter[Obs]) SendMessage(to d_user.ChatID, message d_message.Message, platform string) error {
	m.messages = append(m.messages, message)
	return nil
}

func (m *mockRouter[Obs]) SetObservation(chatID d_user.ChatID, observation string) error {
	m.observation = observation
	m.observationChanged = true
	return nil
}

func (m *mockRouter[Obs]) SetRoute(chatID d_user.ChatID, route string) error {
	m.nextRoute = route
	return nil
}

func (m *mockRouter[Obs]) EndSession(chatID d_user.ChatID, actionId string) error {
	m.endSession = &d_action.EndAction{ID: actionId}
	return nil
}

func (m *mockRouter[Obs]) TransferToMenu(chatID d_user.ChatID, transfer d_action.TransferToMenu, message d_message.Message) error {
	m.transfer = &transfer
	return nil
}

func (m *mockRouter[Obs]) UploadFile(filepath string) (*d_file.File, error) {
	return nil, nil
}

func (m *mockRouter[Obs]) GetFile(fileID string) (*d_file.File, error) {
	return nil, nil
}

// Execute processes a message and returns the result without side effects.
// This is the core method for testing route handlers in isolation.
//
// It handles:
//   - Trigger matching
//   - Loop detection
//   - Handler execution with timeout
//   - Result capture (messages, route changes, etc.)
func (e *Engine[Obs]) Execute(
	userState d_user.UserState[Obs],
	message d_message.Message,
) ExecuteResult[Obs] {
	result := ExecuteResult[Obs]{
		Messages:    []d_message.Message{},
		Observation: userState.Observation,
	}

	// Check for triggers
	preRoute := e.applyTriggers(message.EntireText())
	route := userState.Route

	if preRoute != "" && preRoute != route.Current() {
		log.Printf("[INFO] Triggered route change to: %s", preRoute)
		result.Redirect = &d_action.RedirectResponse{TargetRoute: preRoute}
		return result
	}

	// Check for loops
	loopLimit := e.defaultOptions.LoopCount.Count
	repeated := route.CurrentRepeated()
	if repeated > loopLimit && route.Current() != e.defaultOptions.LoopCount.Route {
		log.Printf("[ERROR] Loop detected for route: %s", route.Current())
		result.Redirect = &d_action.RedirectResponse{
			TargetRoute: e.defaultOptions.LoopCount.Route,
		}
		return result
	}

	// Get route handler
	routeFunc, exists := e.routes[route.Current()]
	if !exists {
		result.Error = fmt.Errorf("route not found: %s", route.Current())
		return result
	}

	// Create mock router to capture outputs
	mock := &mockRouter[Obs]{}

	// Create context with mock router
	ctx, cancel := d_context.NewChatContext(
		userState,
		message,
		mock,
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
	case handlerResult := <-resultChan:
		result.Messages = mock.messages
		result.ObservationChanged = mock.observationChanged

		if handlerResult == nil {
			result.NextRoute = userState.Route.Current()
			return result
		}

		// Process handler result
		switch r := handlerResult.(type) {
		case *d_action.RedirectResponse:
			result.Redirect = r
		case d_action.RedirectResponse:
			result.Redirect = &r
		case *d_action.EndAction:
			result.EndSession = r
		case d_action.EndAction:
			result.EndSession = &r
		case *d_action.TransferToMenu:
			result.Transfer = r
		case d_action.TransferToMenu:
			result.Transfer = &r
		default:
			// Assume it's a route
			if rr, ok := handlerResult.(route_return.RouteReturn); ok {
				if route, isRoute := rr.(*d_route.Route); isRoute {
					result.NextRoute = route.Current()
				} else if route, isRoute := rr.(d_route.Route); isRoute {
					result.NextRoute = route.Current()
				}
			}
		}

		return result

	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("[ERROR] Handler timeout for route: %s", route.Current())
			result.TimedOut = true
			result.Redirect = &d_action.RedirectResponse{
				TargetRoute: routeFunc.HandlerOptions.Timeout.Route,
			}
		} else {
			result.Error = ctx.Err()
		}
		return result
	}
}

// ExecuteWithRouter processes a message using a real router for side effects.
// This is used by the App to execute handlers with actual I/O.
func (e *Engine[Obs]) ExecuteWithRouter(
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

	// Create context with real router
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
