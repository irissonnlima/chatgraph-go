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

type TimeoutRouteOps struct {
	// Duration specifies the maximum duration allowed for the handler execution.
	Duration time.Duration
	Route    string
}

type LoopCountRouteOps struct {
	Count int
	Route string
}

// RouterHandlerOptions configures the behavior and constraints for router handler execution.
// It provides settings for error tracking, execution time limits, and route protection
// to ensure robust and controlled request processing.
type RouterHandlerOptions struct {
	LoopCount *LoopCountRouteOps

	// Timeout specifies the maximum duration allowed for the handler execution.
	// Defaults to DEFAULT_TIMEOUT if not specified.
	Timeout *TimeoutRouteOps

	// Protected configures route protection settings.
	// When enabled, unauthorized users will be redirected to the specified route.
	Protected *ProtectedRouteOps

	Triggers []RouteTrigger
}

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
