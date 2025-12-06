// Package service provides the main chatbot application service.
// It handles route registration, message processing, and result handling.
package service

import (
	"log"

	route_return "github.com/irissonnlima/chatgraph-go/core/domain"
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_route "github.com/irissonnlima/chatgraph-go/core/domain/route"
	d_router "github.com/irissonnlima/chatgraph-go/core/domain/router"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
	adapter_input "github.com/irissonnlima/chatgraph-go/core/ports/adapters/input"
	adapter_output "github.com/irissonnlima/chatgraph-go/core/ports/adapters/output"
)

// ChatbotApp is the main application structure that manages routes and message processing.
// It is generic over Obs, allowing custom observation data types to be used throughout
// the application.
type ChatbotApp[Obs any] struct {
	// engine handles route registration and execution logic.
	engine *Engine[Obs]
	// messageReceiver handles message queue consumption.
	messageReceiver adapter_input.IMessageReceiver[Obs]
	// routerActions provides messaging and session management capabilities.
	routerActions adapter_output.IBotExecutor
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
	routerActions adapter_output.IBotExecutor,
	defaultOptions ...d_router.RouterHandlerOptions,
) *ChatbotApp[Obs] {
	return &ChatbotApp[Obs]{
		engine:          NewEngine[Obs](defaultOptions...),
		messageReceiver: queueAdapter,
		routerActions:   routerActions,
	}
}

// RegisterRoute registers a route handler with optional configuration.
// The route name is case-sensitive and must be unique.
func (app *ChatbotApp[Obs]) RegisterRoute(
	route string,
	handler d_router.RouteHandler[Obs],
	options ...d_router.RouterHandlerOptions,
) {
	app.engine.RegisterRoute(route, handler, options...)
}

// RegisterTrigger registers a global trigger that applies to all routes.
// Triggers are regex patterns that, when matched, redirect to a specific route.
func (app *ChatbotApp[Obs]) RegisterTrigger(trigger d_router.RouteTrigger) {
	app.engine.RegisterTrigger(trigger)
}

// GetEngine returns the underlying engine for testing purposes.
func (app *ChatbotApp[Obs]) GetEngine() *Engine[Obs] {
	return app.engine
}

// HandleMessage processes an incoming message by finding and executing
// the appropriate route handler based on the user's current route.
// Returns an error if no handler is registered for the route or if
// the handler returns an error.
func (app *ChatbotApp[Obs]) HandleMessage(userState d_user.UserState[Obs], message d_message.Message) (err error) {
	result, err := app.engine.ExecuteWithRouter(userState, message, app.routerActions)
	if err != nil {
		return err
	}

	if result == nil {
		// If handler returns nil, stay on the same route
		nextRoute := userState.Route.Next(userState.Route.Current())
		app.handleResult(userState, message, nextRoute)
		return nil
	}

	// Check if result is a redirect (from triggers or loops)
	if redirect, ok := result.(*d_action.RedirectResponse); ok {
		return app.handleRedirect(userState, message, *redirect)
	}
	if redirect, ok := result.(d_action.RedirectResponse); ok {
		return app.handleRedirect(userState, message, redirect)
	}

	app.handleResult(userState, message, result)
	return nil
}

// handleRedirect processes a redirect action by executing the target route.
func (app *ChatbotApp[Obs]) handleRedirect(
	userState d_user.UserState[Obs],
	message d_message.Message,
	redirect d_action.RedirectResponse,
) error {
	// Update user state with new route
	userState.Route = userState.Route.Next(redirect.TargetRoute)
	return app.HandleMessage(userState, message)
}

// handleResult processes the result of a route handler.
func (app *ChatbotApp[Obs]) handleResult(
	userState d_user.UserState[Obs],
	message d_message.Message,
	result route_return.RouteReturn,
) {
	chatID := userState.ChatID

	switch r := result.(type) {
	case *d_action.EndAction:
		app.routerActions.EndSession(chatID, r.ID)

	case d_action.EndAction:
		app.routerActions.EndSession(chatID, r.ID)

	case *d_action.RedirectResponse:
		app.handleRedirect(userState, message, *r)

	case d_action.RedirectResponse:
		app.handleRedirect(userState, message, r)

	case *d_action.TransferToMenu:
		app.routerActions.TransferToMenu(chatID, *r, message)

	case d_action.TransferToMenu:
		app.routerActions.TransferToMenu(chatID, r, message)

	default:
		// Assume it's a route
		switch route := result.(type) {
		case *d_route.Route:
			app.routerActions.SetRoute(chatID, route.Current())
		case d_route.Route:
			app.routerActions.SetRoute(chatID, route.Current())
		}
	}
}

// checkHealthRoutes validates the registered routes before starting the application.
func (app *ChatbotApp[Obs]) checkHealthRoutes() error {
	return app.engine.ValidateRoutes()
}

// Start begins consuming messages from the message receiver in an infinite loop.
// It only returns an error in case of a critical failure.
// Non-critical errors from HandleMessage are logged but do not stop the consumer.
func (app *ChatbotApp[Obs]) Start() error {
	if err := app.checkHealthRoutes(); err != nil {
		log.Printf("[ERROR] Failed to setup routes: %v", err)
		return err
	}

	messages := app.messageReceiver.ConsumeMessage()

	for msg := range messages {
		err := app.HandleMessage(msg.UserState, msg.Message)
		if err != nil {
			log.Printf("[ERROR] Failed to handle message: %v", err)
		}
	}

	log.Println("[CRITICAL] Message channel closed unexpectedly")
	return nil
}
