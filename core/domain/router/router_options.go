package d_router

import (
	"time"
)

// Default values for route handler options.
var (
	// DEFAULT_TIMEOUT is the default maximum duration allowed for route handler execution.
	DEFAULT_TIMEOUT = TimeoutRouteOps{
		Duration: 5 * time.Minute,
		Route:    "timeout_route",
	}

	// DEFAULT_LOOP_COUNT is the default number of loops allowed before stopping handler execution.
	DEFAULT_LOOP_COUNT = LoopCountRouteOps{
		Count: 3,
		Route: "loop_route",
	}
)

// ProtectedRouteOps configures route protection settings.
// When enabled, users who don't meet the protection criteria will be redirected
// to the specified route instead of accessing the protected handler.
type ProtectedRouteOps struct {
	// Route is the route name to redirect to if the user is not allowed access.
	Route string
}

// TimeoutRouteOps configures timeout behavior for route handler execution.
// When a handler exceeds the specified duration, the user is redirected
// to the specified route and the handler execution is cancelled.
type TimeoutRouteOps struct {
	// Duration specifies the maximum duration allowed for the handler execution.
	// If the handler does not complete within this time, it will be cancelled.
	Duration time.Duration
	// Route is the route name to redirect to when a timeout occurs.
	Route string
}

// LoopCountRouteOps configures loop protection for route handler execution.
// This prevents infinite redirect loops by limiting the number of consecutive
// route executions and redirecting to a fallback route when exceeded.
type LoopCountRouteOps struct {
	// Count is the maximum number of consecutive route executions allowed.
	// When this limit is exceeded, the user is redirected to the fallback Route.
	Count int
	// Route is the route name to redirect to when the loop limit is exceeded.
	Route string
}

// RouterHandlerOptions configures the behavior and constraints for router handler execution.
// It provides settings for error tracking, execution time limits, and route protection
// to ensure robust and controlled request processing.
type RouterHandlerOptions struct {
	// LoopCount configures loop protection to prevent infinite redirect loops.
	// If nil, DEFAULT_LOOP_COUNT is used.
	LoopCount *LoopCountRouteOps

	// Timeout specifies the maximum duration allowed for the handler execution.
	// Defaults to DEFAULT_TIMEOUT if not specified.
	Timeout *TimeoutRouteOps

	// Protected configures route protection settings.
	// When enabled, unauthorized users will be redirected to the specified route.
	Protected *ProtectedRouteOps

	// Triggers is a list of regex-based triggers that can automatically redirect
	// the conversation to a different route based on message content.
	Triggers []RouteTrigger
}

// SetOps merges the options from another RouterHandlerOptions into this one.
// Non-nil values in the other options will override the corresponding values in this instance.
// This allows for selective option overriding while preserving defaults.
func (o *RouterHandlerOptions) SetOps(other RouterHandlerOptions) {

	if other.LoopCount != nil {
		o.LoopCount = other.LoopCount
	}
	if other.Timeout != nil {
		o.Timeout = other.Timeout
	}
	if other.Protected != nil {
		o.Protected = other.Protected
	}
	if len(other.Triggers) > 0 {
		o.Triggers = other.Triggers
	}
}

func (o *RouterHandlerOptions) GetRhoRoutes() []string {
	rhoRoutes := []string{}

	if o.Timeout != nil {
		rhoRoutes = append(rhoRoutes, o.Timeout.Route)
	}
	if o.LoopCount != nil {
		rhoRoutes = append(rhoRoutes, o.LoopCount.Route)
	}
	if o.Protected != nil {
		rhoRoutes = append(rhoRoutes, o.Protected.Route)
	}

	return rhoRoutes
}
