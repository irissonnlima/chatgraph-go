// Package d_guard provides the guard function type and default implementation
// for route-level authorization in chatgraph.
package d_guard

import (
	route_return "github.com/irissonnlima/chatgraph-go/core/domain"
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

// GuardFunc is a function that checks whether a user has sufficient access
// for a given auth level. It returns nil if access is granted, or a RouteReturn
// (e.g., RedirectResponse, TransferToMenu) to deny access and redirect.
type GuardFunc[Obs any] func(userState d_user.UserState[Obs], authLevel d_user.AuthLevel) route_return.RouteReturn

// DefaultGuard implements the standard authorization check.
// It verifies user.Identity.AuthLevel against the required level,
// and for "internal" checks whether user.Internal is populated.
// When access is denied, it redirects to the specified deniedRoute.
func DefaultGuard[Obs any](deniedRoute string) GuardFunc[Obs] {
	return func(userState d_user.UserState[Obs], authLevel d_user.AuthLevel) route_return.RouteReturn {
		if !userState.User.Identity.AuthLevel.IsAtLeast(authLevel) {
			return &d_action.RedirectResponse{TargetRoute: deniedRoute}
		}
		return nil
	}
}

// InternalGuard checks that the user has internal HR data (employee relationship).
// When access is denied, it redirects to the specified deniedRoute.
func InternalGuard[Obs any](deniedRoute string) GuardFunc[Obs] {
	return func(userState d_user.UserState[Obs], _ d_user.AuthLevel) route_return.RouteReturn {
		if userState.User.Internal == nil {
			return &d_action.RedirectResponse{TargetRoute: deniedRoute}
		}
		return nil
	}
}
