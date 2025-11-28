// Package d_action provides action types that can be returned from route handlers.
// These actions determine what happens after a route handler completes execution.
package d_action

// EndAction represents the termination of a conversation session.
// When returned from a route handler, it signals that the conversation
// should be ended and any cleanup actions should be performed.
type EndAction struct {
	// ID is the unique identifier for this end action.
	ID string
	// Name is a descriptive name for the end action.
	Name string
	// DepartmentID is the department associated with this action.
	DepartmentID int
	// Observation contains any additional notes about the action.
	Observation string
	// LastUpdate is the timestamp of the last update.
	LastUpdate string
}

// IsRouteReturn implements the RouteReturn interface.
func (EndAction) IsRouteReturn() {}

// RedirectResponse indicates an immediate redirection to another route.
// Unlike setting the next route, a redirect will immediately execute
// the handler for the target route without waiting for user input.
type RedirectResponse struct {
	// TargetRoute is the route to redirect to.
	TargetRoute string
}

// IsRouteReturn implements the RouteReturn interface.
func (RedirectResponse) IsRouteReturn() {}

// TransferToMenu indicates a transfer of the conversation to a different menu.
type TransferToMenu struct {
	// MenuID is the identifier of the menu to transfer to.
	MenuID int
	Route  string
}

// IsRouteReturn implements the RouteReturn interface.
func (TransferToMenu) IsRouteReturn() {}
