// Package service provides the main chatbot application service.
// It handles route registration, message processing, and result handling.
package service

import (
	"context"
	"log"
	"regexp"
	"strings"

	route_return "chatgraph/core/domain"
	d_action "chatgraph/core/domain/action"
	d_context "chatgraph/core/domain/context"
	d_message "chatgraph/core/domain/message"
	d_route "chatgraph/core/domain/route"
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

// NewChatbotApp creates a new ChatbotApp instance with the provided adapters.
// The queueAdapter is used to consume incoming messages from the message broker.
// The routerActions provides messaging and session management capabilities.
// Optional defaultOptions can be provided to set default handler options.
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

// RegisterRoute registers a route handler for the specified route name.
// When a message is received and the user's current route matches the
// registered route, the corresponding handler will be invoked.
// The route name is normalized to lowercase and trimmed of whitespace.
// Optional RouterHandlerOptions can be provided to customize handler behavior;
// if not provided, default options are used.
// Panics if a route with the same name is already registered.
func (app *ChatbotApp[Obs]) RegisterRoute(
	route string,
	handler d_router.RouteHandler[Obs],
	options ...d_router.RouterHandlerOptions,
) {

	route = strings.TrimSpace(route)
	route = strings.ToLower(route)

	if _, exists := app.routes[route]; exists {
		panic("route already registered: " + route)
	}

	opts := app.defaultOptions
	if len(options) > 0 {
		d := options[0]
		opts.SetOps(d)
	}

	app.routes[route] = d_router.RouterHandlerAdmnistrator[Obs]{
		HandlerOptions: opts,
		Handler:        handler,
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

	if route.CurrentRepeated() > loopLimit {
		log.Printf("[ERROR] Loop detected for route: %s", route.Current())
		redirect := d_action.RedirectResponse{
			TargetRoute: app.defaultOptions.Timeout.Route,
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

// handleResult processes the route handler result based on its type.
// It dispatches to the appropriate handler method for each result type.
func (app *ChatbotApp[Obs]) handleResult(userState d_user.UserState[Obs], message d_message.Message, result route_return.RouteReturn) {
	switch r := result.(type) {
	case d_action.EndAction:
		err := app.handleEndAction(userState, r)
		if err != nil {
			log.Printf("[ERROR] failed to end session: %v", err)
		}
		return
	case d_route.Route:
		err := app.handleNextRoute(userState, r)
		if err != nil {
			log.Printf("[ERROR] failed to handle next route: %v", err)
		}
		return
	case d_action.TransferToMenu:
		app.handleTransferToMenu(userState, message, r)
		return
	case d_action.RedirectResponse:
		err := app.handleRedirect(userState, message, r)
		if err != nil {
			log.Printf("[ERROR] failed to handle redirect: %v", err)
		}
		return
	default:
		log.Printf("[ERROR] unknown result type: %T", result)
		return
	}
}

// handleEndAction processes an EndAction result.
// This terminates the conversation session and performs any necessary cleanup.
func (app *ChatbotApp[Obs]) handleEndAction(userState d_user.UserState[Obs], action d_action.EndAction) error {
	return app.routerActions.EndSession(userState.ChatID, action.ID)
}

// handleNextRoute processes a Route result.
// It sets the next route for the user but does not execute the handler.
// The handler will be executed when the user sends their next message.
func (app *ChatbotApp[Obs]) handleNextRoute(userState d_user.UserState[Obs], next d_route.Route) error {

	return app.routerActions.SetRoute(userState.ChatID, next.Current())
}

func (app *ChatbotApp[Obs]) handleRedirect(userState d_user.UserState[Obs], message d_message.Message, redirect d_action.RedirectResponse) error {
	err := app.routerActions.SetRoute(userState.ChatID, redirect.TargetRoute)
	if err != nil {
		return err
	}

	// Immediately execute the target route handler
	newUserState := userState
	newUserState.Route = newUserState.Route.Next(redirect.TargetRoute)

	return app.HandleMessage(newUserState, message)
}

func (app *ChatbotApp[Obs]) handleTransferToMenu(userState d_user.UserState[Obs], message d_message.Message, transfer d_action.TransferToMenu) error {
	return app.routerActions.TransferToMenu(userState.ChatID, transfer, message)
}
