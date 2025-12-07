package service

import (
	"testing"
	"time"

	route_return "github.com/irissonnlima/chatgraph-go/core/domain"
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_context "github.com/irissonnlima/chatgraph-go/core/domain/context"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_route "github.com/irissonnlima/chatgraph-go/core/domain/route"
	d_router "github.com/irissonnlima/chatgraph-go/core/domain/router"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

type TestObs struct {
	Value string
}

// TestNewEngine_DefaultOptions tests NewEngine with default options.
func TestNewEngine_DefaultOptions(t *testing.T) {
	engine := NewEngine[TestObs]()

	if engine == nil {
		t.Fatal("NewEngine returned nil")
	}

	if engine.routes == nil {
		t.Error("routes map is nil")
	}

	if engine.defaultOptions.Timeout == nil {
		t.Error("default timeout is nil")
	}

	if engine.defaultOptions.LoopCount == nil {
		t.Error("default loop count is nil")
	}
}

// TestNewEngine_CustomOptions tests NewEngine with custom options.
func TestNewEngine_CustomOptions(t *testing.T) {
	customTimeout := &d_router.TimeoutRouteOps{
		Duration: 10 * time.Minute,
		Route:    "custom_timeout",
	}
	customLoop := &d_router.LoopCountRouteOps{
		Count: 5,
		Route: "custom_loop",
	}

	engine := NewEngine[TestObs](d_router.RouterHandlerOptions{
		Timeout:   customTimeout,
		LoopCount: customLoop,
	})

	if engine.defaultOptions.Timeout.Duration != 10*time.Minute {
		t.Errorf("expected timeout duration 10m, got %v", engine.defaultOptions.Timeout.Duration)
	}

	if engine.defaultOptions.Timeout.Route != "custom_timeout" {
		t.Errorf("expected timeout route 'custom_timeout', got %s", engine.defaultOptions.Timeout.Route)
	}

	if engine.defaultOptions.LoopCount.Count != 5 {
		t.Errorf("expected loop count 5, got %d", engine.defaultOptions.LoopCount.Count)
	}
}

// TestRegisterRoute tests route registration.
func TestRegisterRoute(t *testing.T) {
	engine := NewEngine[TestObs]()

	handler := func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	}

	engine.RegisterRoute("test_route", handler)

	if _, exists := engine.routes["test_route"]; !exists {
		t.Error("route 'test_route' was not registered")
	}
}

// TestRegisterRoute_WithOptions tests route registration with custom options.
func TestRegisterRoute_WithOptions(t *testing.T) {
	engine := NewEngine[TestObs]()

	handler := func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	}

	customTimeout := &d_router.TimeoutRouteOps{
		Duration: 30 * time.Second,
		Route:    "route_timeout",
	}

	engine.RegisterRoute("test_route", handler, d_router.RouterHandlerOptions{
		Timeout: customTimeout,
	})

	route := engine.routes["test_route"]
	if route.HandlerOptions.Timeout.Duration != 30*time.Second {
		t.Errorf("expected route timeout 30s, got %v", route.HandlerOptions.Timeout.Duration)
	}
}

// TestRegisterTrigger tests trigger registration.
func TestRegisterTrigger(t *testing.T) {
	engine := NewEngine[TestObs]()

	trigger := d_router.RouteTrigger{
		Regex: "^help$",
		Route: "help_route",
	}

	engine.RegisterTrigger(trigger)

	if len(engine.routeTriggers) != 1 {
		t.Errorf("expected 1 trigger, got %d", len(engine.routeTriggers))
	}

	if engine.routeTriggers[0].Route != "help_route" {
		t.Errorf("expected trigger route 'help_route', got %s", engine.routeTriggers[0].Route)
	}
}

// TestApplyTriggers_Match tests trigger matching.
func TestApplyTriggers_Match(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterTrigger(d_router.RouteTrigger{
		Regex: "^help$",
		Route: "help_route",
	})

	result := engine.applyTriggers("help")
	if result != "help_route" {
		t.Errorf("expected 'help_route', got '%s'", result)
	}
}

// TestApplyTriggers_NoMatch tests trigger not matching.
func TestApplyTriggers_NoMatch(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterTrigger(d_router.RouteTrigger{
		Regex: "^help$",
		Route: "help_route",
	})

	result := engine.applyTriggers("hello")
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

// TestApplyTriggers_InvalidRegex tests handling of invalid regex.
func TestApplyTriggers_InvalidRegex(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterTrigger(d_router.RouteTrigger{
		Regex: "[invalid",
		Route: "some_route",
	})

	result := engine.applyTriggers("test")
	if result != "" {
		t.Errorf("expected empty string for invalid regex, got '%s'", result)
	}
}

// TestExecute_Success tests successful handler execution.
func TestExecute_Success(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterRoute("start", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		ctx.SendTextMessage("Hello!")
		return ctx.NextRoute("next")
	})

	userState := d_user.UserState[TestObs]{
		Route: d_route.Route{
			History:   []string{"start"},
			Separator: '/',
		},
	}

	msg := d_message.Message{}
	mock := newMockExecutor()

	result, err := engine.Execute(userState, msg, mock)

	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	route, ok := result.(d_route.Route)
	if !ok {
		t.Fatalf("expected Route, got %T", result)
	}

	if route.Current() != "next" {
		t.Errorf("expected current route 'next', got '%s'", route.Current())
	}
}

// TestExecute_RouteNotFound tests execution with non-existent route.
func TestExecute_RouteNotFound(t *testing.T) {
	engine := NewEngine[TestObs]()

	userState := d_user.UserState[TestObs]{
		Route: d_route.Route{
			History:   []string{"nonexistent"},
			Separator: '/',
		},
	}

	msg := d_message.Message{}
	mock := newMockExecutor()

	_, err := engine.Execute(userState, msg, mock)

	if err == nil {
		t.Error("expected error for non-existent route")
	}
}

// TestExecute_TriggerMatch tests execution when trigger matches.
func TestExecute_TriggerMatch(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterRoute("start", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})

	engine.RegisterTrigger(d_router.RouteTrigger{
		Regex: "^help$",
		Route: "help_route",
	})

	userState := d_user.UserState[TestObs]{
		Route: d_route.Route{
			History:   []string{"start"},
			Separator: '/',
		},
	}

	msg := d_message.Message{
		TextMessage: d_message.TextMessage{Detail: "help"},
	}
	mock := newMockExecutor()

	result, err := engine.Execute(userState, msg, mock)

	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	redirect, ok := result.(*d_action.RedirectResponse)
	if !ok {
		t.Fatalf("expected RedirectResponse, got %T", result)
	}

	if redirect.TargetRoute != "help_route" {
		t.Errorf("expected target route 'help_route', got '%s'", redirect.TargetRoute)
	}
}

// TestExecute_LoopDetection tests loop detection.
func TestExecute_LoopDetection(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterRoute("loop_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})

	engine.RegisterRoute("test", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})

	// Simulate a route that has been repeated more than the loop limit
	userState := d_user.UserState[TestObs]{
		Route: d_route.Route{
			History:   []string{"test", "test", "test", "test", "test"}, // 5 times, exceeds default limit of 3
			Separator: '/',
		},
	}

	msg := d_message.Message{}
	mock := newMockExecutor()

	result, err := engine.Execute(userState, msg, mock)

	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	redirect, ok := result.(*d_action.RedirectResponse)
	if !ok {
		t.Fatalf("expected RedirectResponse for loop detection, got %T", result)
	}

	if redirect.TargetRoute != "loop_route" {
		t.Errorf("expected redirect to 'loop_route', got '%s'", redirect.TargetRoute)
	}
}

// TestExecute_NilReturn tests handler returning nil.
func TestExecute_NilReturn(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterRoute("start", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})

	userState := d_user.UserState[TestObs]{
		Route: d_route.Route{
			History:   []string{"start"},
			Separator: '/',
		},
	}

	msg := d_message.Message{}
	mock := newMockExecutor()

	result, err := engine.Execute(userState, msg, mock)

	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	route, ok := result.(d_route.Route)
	if !ok {
		t.Fatalf("expected Route when handler returns nil, got %T", result)
	}

	// When handler returns nil, it should stay on current route
	if route.Current() != "start" {
		t.Errorf("expected current route 'start', got '%s'", route.Current())
	}
}

// TestExecute_Timeout tests handler timeout.
func TestExecute_Timeout(t *testing.T) {
	customTimeout := &d_router.TimeoutRouteOps{
		Duration: 50 * time.Millisecond,
		Route:    "timeout_route",
	}

	engine := NewEngine[TestObs](d_router.RouterHandlerOptions{
		Timeout: customTimeout,
	})

	engine.RegisterRoute("slow_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	engine.RegisterRoute("timeout_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})

	userState := d_user.UserState[TestObs]{
		Route: d_route.Route{
			History:   []string{"slow_route"},
			Separator: '/',
		},
	}

	msg := d_message.Message{}
	mock := newMockExecutor()

	result, err := engine.Execute(userState, msg, mock)

	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	redirect, ok := result.(*d_action.RedirectResponse)
	if !ok {
		t.Fatalf("expected RedirectResponse for timeout, got %T", result)
	}

	if redirect.TargetRoute != "timeout_route" {
		t.Errorf("expected redirect to 'timeout_route', got '%s'", redirect.TargetRoute)
	}
}

// TestValidateRoutes_Success tests successful route validation.
func TestValidateRoutes_Success(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterRoute("start", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})
	engine.RegisterRoute("timeout_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})
	engine.RegisterRoute("loop_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})

	err := engine.ValidateRoutes()

	if err != nil {
		t.Errorf("ValidateRoutes returned error: %v", err)
	}
}

// TestValidateRoutes_MissingStart tests validation without start route.
func TestValidateRoutes_MissingStart(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterRoute("other", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})

	err := engine.ValidateRoutes()

	if err == nil {
		t.Error("expected error for missing 'start' route")
	}
}

// TestValidateRoutes_MissingTriggerRoute tests validation with missing trigger route.
func TestValidateRoutes_MissingTriggerRoute(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterRoute("start", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})
	engine.RegisterRoute("timeout_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})
	engine.RegisterRoute("loop_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})

	engine.RegisterTrigger(d_router.RouteTrigger{
		Regex: "^help$",
		Route: "help_route", // This route is not registered
	})

	err := engine.ValidateRoutes()

	if err == nil {
		t.Error("expected error for missing trigger route")
	}
}

// TestValidateRoutes_MissingDefaultRoute tests validation with missing default option route.
func TestValidateRoutes_MissingDefaultRoute(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterRoute("start", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})
	// Missing timeout_route and loop_route

	err := engine.ValidateRoutes()

	if err == nil {
		t.Error("expected error for missing default option routes")
	}
}

// TestValidateRoutes_MissingRouteHandlerTrigger tests validation with missing handler trigger route.
func TestValidateRoutes_MissingRouteHandlerTrigger(t *testing.T) {
	engine := NewEngine[TestObs]()

	engine.RegisterRoute("start", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	}, d_router.RouterHandlerOptions{
		Triggers: []d_router.RouteTrigger{
			{Regex: "^test$", Route: "missing_trigger_route"},
		},
	})
	engine.RegisterRoute("timeout_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})
	engine.RegisterRoute("loop_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return nil
	})

	err := engine.ValidateRoutes()

	if err == nil {
		t.Error("expected error for missing handler trigger route")
	}
}

// TestExecute_TriggerSameRoute tests that trigger doesn't redirect when already on that route.
func TestExecute_TriggerSameRoute(t *testing.T) {
	engine := NewEngine[TestObs]()

	executed := false
	engine.RegisterRoute("help_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		executed = true
		return nil
	})

	engine.RegisterTrigger(d_router.RouteTrigger{
		Regex: "^help$",
		Route: "help_route",
	})

	userState := d_user.UserState[TestObs]{
		Route: d_route.Route{
			History:   []string{"help_route"}, // Already on help_route
			Separator: '/',
		},
	}

	msg := d_message.Message{
		TextMessage: d_message.TextMessage{Detail: "help"},
	}
	mock := newMockExecutor()

	_, err := engine.Execute(userState, msg, mock)

	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !executed {
		t.Error("handler should have been executed when already on trigger route")
	}
}

// TestExecute_LoopOnLoopRoute tests that loop detection doesn't trigger when on loop_route.
func TestExecute_LoopOnLoopRoute(t *testing.T) {
	engine := NewEngine[TestObs]()

	executed := false
	engine.RegisterRoute("loop_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		executed = true
		return nil
	})

	// Simulate being on loop_route multiple times
	userState := d_user.UserState[TestObs]{
		Route: d_route.Route{
			History:   []string{"loop_route", "loop_route", "loop_route", "loop_route", "loop_route"},
			Separator: '/',
		},
	}

	msg := d_message.Message{}
	mock := newMockExecutor()

	_, err := engine.Execute(userState, msg, mock)

	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !executed {
		t.Error("handler should have been executed when on loop_route (no infinite redirect)")
	}
}
