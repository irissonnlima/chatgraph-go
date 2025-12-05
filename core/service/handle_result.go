package service

import (
	route_return "chatgraph/core/domain"
	d_action "chatgraph/core/domain/action"
	d_message "chatgraph/core/domain/message"
	d_route "chatgraph/core/domain/route"
	d_user "chatgraph/core/domain/user"
	"log"
)

// handleEndSession processes an EndAction result.
// This terminates the conversation session and performs any necessary cleanup.
func (app *ChatbotApp[Obs]) handleEndSession(userState d_user.UserState[Obs], action d_action.EndAction) error {
	return app.routerActions.EndSession(userState.ChatID, action.ID)
}

// handleNextRoute processes a Route result.
// It sets the next route for the user but does not execute the handler.
// The handler will be executed when the user sends their next message.
func (app *ChatbotApp[Obs]) handleNextRoute(userState d_user.UserState[Obs], next d_route.Route) error {

	return app.routerActions.SetRoute(userState.ChatID, next.Current())
}

// handleRedirect processes a RedirectResponse result.
// It immediately sets the new route and recursively calls HandleMessage
// to execute the target route handler without waiting for user input.
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

// handleTransferToMenu processes a TransferToMenu result.
// It delegates to the router service to transfer the user to a different menu.
func (app *ChatbotApp[Obs]) handleTransferToMenu(userState d_user.UserState[Obs], message d_message.Message, transfer d_action.TransferToMenu) error {
	return app.routerActions.TransferToMenu(userState.ChatID, transfer, message)
}

// handleResult processes the route handler result based on its type.
// It dispatches to the appropriate handler method for each result type.
func (app *ChatbotApp[Obs]) handleResult(userState d_user.UserState[Obs], message d_message.Message, result route_return.RouteReturn) {
	switch r := result.(type) {
	case d_action.EndAction:
		err := app.handleEndSession(userState, r)
		if err != nil {
			log.Printf("[ERROR] failed to end session: %v", err)
		}
		return
	case *d_action.EndAction:
		err := app.handleEndSession(userState, *r)
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
	case *d_route.Route:
		err := app.handleNextRoute(userState, *r)
		if err != nil {
			log.Printf("[ERROR] failed to handle next route: %v", err)
		}
		return
	case d_action.TransferToMenu:
		app.handleTransferToMenu(userState, message, r)
		return
	case *d_action.TransferToMenu:
		app.handleTransferToMenu(userState, message, *r)
		return
	case d_action.RedirectResponse:
		err := app.handleRedirect(userState, message, r)
		if err != nil {
			log.Printf("[ERROR] failed to handle redirect: %v", err)
		}
		return
	case *d_action.RedirectResponse:
		err := app.handleRedirect(userState, message, *r)
		if err != nil {
			log.Printf("[ERROR] failed to handle redirect: %v", err)
		}
		return
	default:
		log.Printf("[ERROR] unknown result type: %T", result)
		return
	}
}
