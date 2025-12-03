// Package d_action provides action types that can be returned from route handlers.
// These actions determine what happens after a route handler completes execution.
package dto_action

// EndAction represents the termination of a conversation session.
// When returned from a route handler, it signals that the conversation
// should be ended and any cleanup actions should be performed.
type EndAction struct {
	// ID is the unique identifier for this end action.
	ID string `json:"id"`
	// Name is a descriptive name for the end action.
	Name string `json:"name,omitempty"`
}

// RedirectResponse indicates an immediate redirection to another route.
// Unlike setting the next route, a redirect will immediately execute
// the handler for the target route without waiting for user input.
type RedirectResponse struct {
	// TargetRoute is the route to redirect to.
	TargetRoute string `json:"target_route"`
}

// TransferToMenu indicates a transfer of the conversation to a different menu.
type TransferToMenu struct {
	// MenuID is the identifier of the menu to transfer to.
	MenuID int    `json:"menu_id"`
	Route  string `json:"route"`
}
