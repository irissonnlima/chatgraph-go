// Package service provides the main chatbot application service.
// It handles route registration, message processing, and result handling.
package service

import (
	"context"
	"log"
	"regexp"

	route_return "chatgraph/core/domain"
	d_action "chatgraph/core/domain/action"
	d_context "chatgraph/core/domain/context"
	d_message "chatgraph/core/domain/message"
	d_router "chatgraph/core/domain/router"
	d_user "chatgraph/core/domain/user"
	adapter_input "chatgraph/core/ports/adapters/input"
	adapter_output "chatgraph/core/ports/adapters/output"
)

// ChatbotApp is the main application structure that manages routes and message processing.
// It is generic over Obs, allowing custom observation data types to be used throughout
// the application.
type ChatbotApp[Obs any] struct {
	// routes maps route names to their handler administrators.
	routes map[string]d_router.RouterHandlerAdmnistrator[Obs]
	// messageReceiver handles message queue consumption.
	messageReceiver adapter_input.IMessageReceiver[Obs]
	// routerActions provides messaging and session management capabilities.
	routerActions adapter_output.RouterService
	// defaultOptions holds the default handler options.
	defaultOptions d_router.RouterHandlerOptions
	// routeTriggers holds the route triggers configuration.
	routeTriggers []d_router.RouteTrigger
}

/*
NewChatbotApp creates a new ChatbotApp instance with the provided adapters.
The queueAdapter is used to consume incoming messages from the message broker.
The routerActions provides messaging and session management capabilities.

Optional defaultOptions can be provided to set default handler options.

If not provided, the following defaults are used:
  - Timeout: 5 minutes (redirects to "timeout_route")
  - Loop Limit: 3 iterations (redirects to "loop_route")
  - Protected: nil (no protection by default)
*/
func NewChatbotApp[Obs any](
	queueAdapter adapter_input.IMessageReceiver[Obs],
	routerActions adapter_output.RouterService,
	defaultOptions ...d_router.RouterHandlerOptions,
) *ChatbotApp[Obs] {

	defaultOpts := d_router.RouterHandlerOptions{
		Timeout:   &d_router.DEFAULT_TIMEOUT,
		LoopCount: &d_router.DEFAULT_LOOP_COUNT,
		Protected: nil,
	}

	if len(defaultOptions) > 0 {
		d := defaultOptions[0]
		defaultOpts.SetOps(d)
	}

	return &ChatbotApp[Obs]{
		routes:          make(map[string]d_router.RouterHandlerAdmnistrator[Obs]),
		messageReceiver: queueAdapter,
		routerActions:   routerActions,
		defaultOptions:  defaultOpts,
	}
}

// applyTriggers checks if the message matches any trigger regex.
// If a match is found, returns the route associated with the trigger.
// Returns empty string if no trigger matches.
func (app *ChatbotApp[Obs]) applyTriggers(messageText string) string {
	for _, trigger := range app.routeTriggers {
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

// HandleMessage processes an incoming message by finding and executing
// the appropriate route handler based on the user's current route.
// Returns an error if no handler is registered for the route or if
// the handler returns an error.
func (app *ChatbotApp[Obs]) HandleMessage(userState d_user.UserState[Obs], message d_message.Message) (err error) {
	// *Pre-process message to check for triggers and loops
	preRoute := app.applyTriggers(message.EntireText())

	route := userState.Route
	if preRoute != "" && preRoute != route.Current() {
		log.Printf("[INFO] Triggered route change to: %s", preRoute)
		return app.handleRedirect(
			userState,
			message,
			d_action.RedirectResponse{TargetRoute: preRoute},
		)
	}

	loopLimit := app.defaultOptions.LoopCount.Count
	repeated := route.CurrentRepeated()
	if repeated > loopLimit && route.Current() != app.defaultOptions.LoopCount.Route {
		log.Printf("[ERROR] Loop detected for route: %s", route.Current())
		redirect := d_action.RedirectResponse{
			TargetRoute: app.defaultOptions.LoopCount.Route,
		}
		return app.handleRedirect(userState, message, redirect)
	}

	// *Execute the route handler
	routeFunc, exists := app.routes[route.Current()]
	if !exists {
		panic("Route not found " + route.Current())
	}

	ctx, cancel := d_context.NewChatContext(
		userState,
		message,
		app.routerActions,
		routeFunc.HandlerOptions.Timeout.Duration,
	)
	defer cancel()

	// Channel to receive the result
	resultChan := make(chan route_return.RouteReturn, 1)

	// Execute handler in goroutine
	go func() {
		resultChan <- routeFunc.Handler(&ctx)
	}()

	// *Post-process message waiting for result or timeout
	select {
	case result := <-resultChan:
		if result == nil {
			// If handler returns nil, stay on the same route
			nextRoute := userState.Route.Next(userState.Route.Current())
			app.handleResult(userState, message, nextRoute)
			return nil
		}
		app.handleResult(userState, message, result)
		return nil

	case <-ctx.Done():
		// Timeout occurred - redirect to timeout route
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("[ERROR] Handler timeout for route: %s", route.Current())
			redirect := d_action.RedirectResponse{
				TargetRoute: routeFunc.HandlerOptions.Timeout.Route,
			}
			return app.handleRedirect(userState, message, redirect)
		}
		return ctx.Err()
	}
}
