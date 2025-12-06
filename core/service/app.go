// Package service provides the main chatbot application service.
// It handles route registration, message processing, and result handling.
package service

import (
	"log"

	route_return "github.com/irissonnlima/chatgraph-go/core/domain"
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_route "github.com/irissonnlima/chatgraph-go/core/domain/route"
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
	// botExecutor provides messaging and session management capabilities.
	botExecutor adapter_output.IBotExecutor
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
	engine *Engine[Obs],
	messageReceiver adapter_input.IMessageReceiver[Obs],
	botExecutor adapter_output.IBotExecutor,
) *ChatbotApp[Obs] {
	return &ChatbotApp[Obs]{
		engine:          engine,
		messageReceiver: messageReceiver,
		botExecutor:     botExecutor,
	}
}

// HandleMessage processes an incoming message by finding and executing
// the appropriate route handler based on the user's current route.
// Returns an error if no handler is registered for the route or if
// the handler returns an error.
func (app *ChatbotApp[Obs]) HandleMessage(userState d_user.UserState[Obs], message d_message.Message) (err error) {
	result, err := app.engine.Execute(userState, message, app.botExecutor)
	if err != nil {
		return err
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
	err := app.botExecutor.SetRoute(userState.ChatID, redirect.TargetRoute)
	if err != nil {
		log.Printf("[ERROR] Failed to set route for chat %v: %v", userState.ChatID, err)
	}
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
	var err error

	switch r := result.(type) {
	case *d_action.EndAction:
		err = app.botExecutor.EndSession(chatID, r.ID)

	case *d_action.RedirectResponse:
		err = app.handleRedirect(userState, message, *r)

	case *d_action.TransferToMenu:
		err = app.botExecutor.TransferToMenu(chatID, *r, message)

	case *d_route.Route:
		err = app.botExecutor.SetRoute(chatID, r.Current())

	case nil:
		err = app.botExecutor.SetRoute(chatID, userState.Route.Current())

	default:
		log.Printf("[WARN] Unhandled route return type for chat %v: %T", chatID, r)
		err = app.botExecutor.SetRoute(chatID, userState.Route.Current())
	}

	if err != nil {
		log.Printf("[ERROR] Failed to handle result for chat %v: %v", chatID, err)
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
