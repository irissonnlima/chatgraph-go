package service

import (
	"testing"

	route_return "github.com/irissonnlima/chatgraph-go/core/domain"
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_context "github.com/irissonnlima/chatgraph-go/core/domain/context"
	d_guard "github.com/irissonnlima/chatgraph-go/core/domain/guard"
	d_logger "github.com/irissonnlima/chatgraph-go/core/domain/logger"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_route "github.com/irissonnlima/chatgraph-go/core/domain/route"
	d_router "github.com/irissonnlima/chatgraph-go/core/domain/router"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

// ==========================================================================
// Guard Integration Tests
// ==========================================================================

func TestExecute_GuardAllowsAccess(t *testing.T) {
	engine := NewEngine[TestObs]()
	engine.SetGuard(d_guard.DefaultGuard[TestObs]("denied_route"))

	engine.RegisterRoute("protected", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		ctx.SendTextMessage("Access granted")
		return ctx.NextRoute("done")
	}, d_router.RouterHandlerOptions{
		Protected: &d_router.ProtectedRouteOps{
			Route:     "denied_route",
			AuthLevel: d_user.AuthRead,
		},
	})

	userState := d_user.UserState[TestObs]{
		ChatID: d_user.ChatID{UserID: "u1", CompanyID: "c1"},
		User: d_user.User{
			Identity: d_user.UserIdentity{AuthLevel: d_user.AuthWrite},
		},
		Route: d_route.Route{
			History:   []string{"protected"},
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
		t.Fatalf("expected Route (access granted), got %T", result)
	}
	if route.Current() != "done" {
		t.Errorf("expected route 'done', got '%s'", route.Current())
	}
}

func TestExecute_GuardDeniesAccess(t *testing.T) {
	engine := NewEngine[TestObs]()
	engine.SetGuard(d_guard.DefaultGuard[TestObs]("denied_route"))

	engine.RegisterRoute("protected", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		t.Error("handler should not be executed when guard denies access")
		return nil
	}, d_router.RouterHandlerOptions{
		Protected: &d_router.ProtectedRouteOps{
			Route:     "denied_route",
			AuthLevel: d_user.AuthWrite,
		},
	})

	userState := d_user.UserState[TestObs]{
		ChatID: d_user.ChatID{UserID: "u1", CompanyID: "c1"},
		User: d_user.User{
			Identity: d_user.UserIdentity{AuthLevel: d_user.AuthRead},
		},
		Route: d_route.Route{
			History:   []string{"protected"},
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
		t.Fatalf("expected RedirectResponse (access denied), got %T", result)
	}
	if redirect.TargetRoute != "denied_route" {
		t.Errorf("expected redirect to 'denied_route', got '%s'", redirect.TargetRoute)
	}
}

func TestExecute_GuardInternalRequired_HasInternal(t *testing.T) {
	engine := NewEngine[TestObs]()
	engine.SetGuard(d_guard.DefaultGuard[TestObs]("denied_route"))

	engine.RegisterRoute("internal_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		ctx.SendTextMessage("Internal access")
		return ctx.NextRoute("done")
	}, d_router.RouterHandlerOptions{
		Protected: &d_router.ProtectedRouteOps{
			Route:    "denied_route",
			Internal: true,
		},
	})

	userState := d_user.UserState[TestObs]{
		ChatID: d_user.ChatID{UserID: "u1", CompanyID: "c1"},
		User: d_user.User{
			Internal: &d_user.UserInternal{
				Matricula: "12345",
				Cargo:     "Dev",
			},
		},
		Route: d_route.Route{
			History:   []string{"internal_route"},
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
		t.Fatalf("expected Route (internal access granted), got %T", result)
	}
	if route.Current() != "done" {
		t.Errorf("expected route 'done', got '%s'", route.Current())
	}
}

func TestExecute_GuardInternalRequired_NoInternal(t *testing.T) {
	engine := NewEngine[TestObs]()
	engine.SetGuard(d_guard.DefaultGuard[TestObs]("denied_route"))

	engine.RegisterRoute("internal_route", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		t.Error("handler should not be executed without internal data")
		return nil
	}, d_router.RouterHandlerOptions{
		Protected: &d_router.ProtectedRouteOps{
			Route:    "denied_route",
			Internal: true,
		},
	})

	userState := d_user.UserState[TestObs]{
		ChatID: d_user.ChatID{UserID: "u1", CompanyID: "c1"},
		User: d_user.User{
			Internal: nil,
		},
		Route: d_route.Route{
			History:   []string{"internal_route"},
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
		t.Fatalf("expected RedirectResponse (internal denied), got %T", result)
	}
	if redirect.TargetRoute != "denied_route" {
		t.Errorf("expected redirect to 'denied_route', got '%s'", redirect.TargetRoute)
	}
}

func TestExecute_NoGuardSet_ProtectedRouteExecutes(t *testing.T) {
	engine := NewEngine[TestObs]()
	// No guard set — protected routes should execute normally

	engine.RegisterRoute("protected", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		return ctx.NextRoute("done")
	}, d_router.RouterHandlerOptions{
		Protected: &d_router.ProtectedRouteOps{
			Route:     "denied_route",
			AuthLevel: d_user.AuthWrite,
		},
	})

	userState := d_user.UserState[TestObs]{
		ChatID: d_user.ChatID{UserID: "u1", CompanyID: "c1"},
		User: d_user.User{
			Identity: d_user.UserIdentity{AuthLevel: d_user.AuthBlocked},
		},
		Route: d_route.Route{
			History:   []string{"protected"},
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
		t.Fatalf("expected Route (no guard = no check), got %T", result)
	}
	if route.Current() != "done" {
		t.Errorf("expected route 'done', got '%s'", route.Current())
	}
}

// ==========================================================================
// Logger Integration Tests
// ==========================================================================

func TestExecute_LoggerAttachedToContext(t *testing.T) {
	dir := t.TempDir()
	manager := d_logger.NewUserLoggerManager(dir, d_logger.LevelDebug)
	defer manager.Close()

	engine := NewEngine[TestObs]()
	engine.SetLoggerManager(manager)

	var loggerReceived *d_logger.UserLogger

	engine.RegisterRoute("start", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		loggerReceived = ctx.Logger
		if ctx.Logger != nil {
			ctx.Logger.Info("test log from handler")
		}
		return ctx.NextRoute("done")
	})

	userState := d_user.UserState[TestObs]{
		ChatID: d_user.ChatID{UserID: "user1", CompanyID: "comp1"},
		Route: d_route.Route{
			History:   []string{"start"},
			Separator: '/',
		},
	}

	msg := d_message.Message{}
	mock := newMockExecutor()

	_, err := engine.Execute(userState, msg, mock)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if loggerReceived == nil {
		t.Error("expected Logger to be attached to context")
	}
}

func TestExecute_NoLoggerManager_LoggerIsNil(t *testing.T) {
	engine := NewEngine[TestObs]()
	// No logger manager set

	var loggerReceived *d_logger.UserLogger

	engine.RegisterRoute("start", func(ctx *d_context.ChatContext[TestObs]) route_return.RouteReturn {
		loggerReceived = ctx.Logger
		return ctx.NextRoute("done")
	})

	userState := d_user.UserState[TestObs]{
		ChatID: d_user.ChatID{UserID: "user1", CompanyID: "comp1"},
		Route: d_route.Route{
			History:   []string{"start"},
			Separator: '/',
		},
	}

	msg := d_message.Message{}
	mock := newMockExecutor()

	_, err := engine.Execute(userState, msg, mock)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if loggerReceived != nil {
		t.Error("expected Logger to be nil when no manager is set")
	}
}
