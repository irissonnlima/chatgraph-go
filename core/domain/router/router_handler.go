package d_router

import (
	route_return "github.com/irissonnlima/chatgraph-go/core/domain"
	d_context "github.com/irissonnlima/chatgraph-go/core/domain/context"
)

// RouteHandler defines the signature of a route handler function.
// Route handlers process incoming messages and return a RouteReturn to indicate
// the next action. They can return:
//   - EndAction: Ends the conversation session
//   - RedirectResponse: Redirects to another route immediately
//   - Route: Sets the next route for the user's next message
//   - TransferToMenu: Transfers the user to a different menu
//
// The Obs type parameter allows custom observation data to be passed through
// the context.
type RouteHandler[Obs any] func(ctx *d_context.ChatContext[Obs]) route_return.RouteReturn

// RouterHandlerAdmnistrator wraps a route handler with its configuration options.
// It associates a RouteHandler with its execution settings like timeout and error tracking.
type RouterHandlerAdmnistrator[Obs any] struct {
	// HandlerOptions contains the configuration for handler execution.
	HandlerOptions RouterHandlerOptions
	// Handler is the route handler function to be executed.
	Handler RouteHandler[Obs]
}
