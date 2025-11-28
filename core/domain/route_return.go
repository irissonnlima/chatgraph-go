// Package route_return defines the interface for route handler return types.
// It provides a common interface that all route handler results must implement.
package route_return

// RouteReturn is a marker interface for all possible route handler results.
// Types implementing this interface can be returned from route handlers
// to indicate the next action to take.
//
// Implementations:
//   - d_action.EndAction: Ends the conversation session
//   - d_action.RedirectResponse: Redirects to another route immediately
//   - d_route.Route: Sets the next route for the user's next message
type RouteReturn interface {
	IsRouteReturn()
}
