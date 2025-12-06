package d_router

import (
	"testing"
	"time"
)

func TestRouterHandlerOptions_SetOps(t *testing.T) {
	t.Run("sets timeout when provided", func(t *testing.T) {
		opts := RouterHandlerOptions{
			Timeout: &DEFAULT_TIMEOUT,
		}

		newTimeout := &TimeoutRouteOps{
			Duration: 10 * time.Second,
			Route:    "custom_timeout",
		}

		opts.SetOps(RouterHandlerOptions{Timeout: newTimeout})

		if opts.Timeout.Duration != 10*time.Second {
			t.Errorf("Timeout.Duration = %v, want 10s", opts.Timeout.Duration)
		}
		if opts.Timeout.Route != "custom_timeout" {
			t.Errorf("Timeout.Route = %v, want custom_timeout", opts.Timeout.Route)
		}
	})

	t.Run("sets loop count when provided", func(t *testing.T) {
		opts := RouterHandlerOptions{
			LoopCount: &DEFAULT_LOOP_COUNT,
		}

		newLoop := &LoopCountRouteOps{
			Count: 10,
			Route: "custom_loop",
		}

		opts.SetOps(RouterHandlerOptions{LoopCount: newLoop})

		if opts.LoopCount.Count != 10 {
			t.Errorf("LoopCount.Count = %v, want 10", opts.LoopCount.Count)
		}
		if opts.LoopCount.Route != "custom_loop" {
			t.Errorf("LoopCount.Route = %v, want custom_loop", opts.LoopCount.Route)
		}
	})

	t.Run("sets protected when provided", func(t *testing.T) {
		opts := RouterHandlerOptions{}

		protected := &ProtectedRouteOps{Route: "login"}

		opts.SetOps(RouterHandlerOptions{Protected: protected})

		if opts.Protected == nil {
			t.Error("Protected should not be nil")
		}
		if opts.Protected.Route != "login" {
			t.Errorf("Protected.Route = %v, want login", opts.Protected.Route)
		}
	})

	t.Run("sets triggers when provided", func(t *testing.T) {
		opts := RouterHandlerOptions{}

		triggers := []RouteTrigger{
			{Regex: "^menu$", Route: "menu"},
		}

		opts.SetOps(RouterHandlerOptions{Triggers: triggers})

		if len(opts.Triggers) != 1 {
			t.Errorf("len(Triggers) = %v, want 1", len(opts.Triggers))
		}
		if opts.Triggers[0].Regex != "^menu$" {
			t.Errorf("Triggers[0].Regex = %v, want ^menu$", opts.Triggers[0].Regex)
		}
	})

	t.Run("does not override when nil", func(t *testing.T) {
		originalTimeout := &TimeoutRouteOps{
			Duration: 5 * time.Minute,
			Route:    "original",
		}
		opts := RouterHandlerOptions{
			Timeout: originalTimeout,
		}

		opts.SetOps(RouterHandlerOptions{Timeout: nil})

		if opts.Timeout != originalTimeout {
			t.Error("Timeout should not have been overridden")
		}
	})
}

func TestRouterHandlerOptions_GetRhoRoutes(t *testing.T) {
	t.Run("returns all routes when all options set", func(t *testing.T) {
		opts := RouterHandlerOptions{
			Timeout:   &TimeoutRouteOps{Route: "timeout"},
			LoopCount: &LoopCountRouteOps{Route: "loop"},
			Protected: &ProtectedRouteOps{Route: "protected"},
		}

		routes := opts.GetRhoRoutes()

		if len(routes) != 3 {
			t.Errorf("len(routes) = %v, want 3", len(routes))
		}

		expected := map[string]bool{"timeout": true, "loop": true, "protected": true}
		for _, r := range routes {
			if !expected[r] {
				t.Errorf("unexpected route: %v", r)
			}
		}
	})

	t.Run("returns empty when no options set", func(t *testing.T) {
		opts := RouterHandlerOptions{}

		routes := opts.GetRhoRoutes()

		if len(routes) != 0 {
			t.Errorf("len(routes) = %v, want 0", len(routes))
		}
	})

	t.Run("returns only timeout when only timeout set", func(t *testing.T) {
		opts := RouterHandlerOptions{
			Timeout: &TimeoutRouteOps{Route: "timeout"},
		}

		routes := opts.GetRhoRoutes()

		if len(routes) != 1 {
			t.Errorf("len(routes) = %v, want 1", len(routes))
		}
		if routes[0] != "timeout" {
			t.Errorf("routes[0] = %v, want timeout", routes[0])
		}
	})
}

func TestDefaultValues(t *testing.T) {
	t.Run("DEFAULT_TIMEOUT has correct values", func(t *testing.T) {
		if DEFAULT_TIMEOUT.Duration != 5*time.Minute {
			t.Errorf("DEFAULT_TIMEOUT.Duration = %v, want 5m", DEFAULT_TIMEOUT.Duration)
		}
		if DEFAULT_TIMEOUT.Route != "timeout_route" {
			t.Errorf("DEFAULT_TIMEOUT.Route = %v, want timeout_route", DEFAULT_TIMEOUT.Route)
		}
	})

	t.Run("DEFAULT_LOOP_COUNT has correct values", func(t *testing.T) {
		if DEFAULT_LOOP_COUNT.Count != 3 {
			t.Errorf("DEFAULT_LOOP_COUNT.Count = %v, want 3", DEFAULT_LOOP_COUNT.Count)
		}
		if DEFAULT_LOOP_COUNT.Route != "loop_route" {
			t.Errorf("DEFAULT_LOOP_COUNT.Route = %v, want loop_route", DEFAULT_LOOP_COUNT.Route)
		}
	})
}
