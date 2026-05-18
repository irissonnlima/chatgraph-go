package d_guard

import (
	"testing"

	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

type testObs struct {
	Value string
}

func TestDefaultGuard_AccessGranted(t *testing.T) {
	guard := DefaultGuard[testObs]("denied_route")

	userState := d_user.UserState[testObs]{
		User: d_user.User{
			Identity: d_user.UserIdentity{AuthLevel: d_user.AuthWrite},
		},
	}

	result := guard(userState, d_user.AuthRead)
	if result != nil {
		t.Errorf("DefaultGuard should return nil when user has sufficient level, got %v", result)
	}
}

func TestDefaultGuard_AccessGranted_EqualLevel(t *testing.T) {
	guard := DefaultGuard[testObs]("denied_route")

	userState := d_user.UserState[testObs]{
		User: d_user.User{
			Identity: d_user.UserIdentity{AuthLevel: d_user.AuthRead},
		},
	}

	result := guard(userState, d_user.AuthRead)
	if result != nil {
		t.Errorf("DefaultGuard should return nil when user has equal level, got %v", result)
	}
}

func TestDefaultGuard_AccessDenied(t *testing.T) {
	guard := DefaultGuard[testObs]("denied_route")

	userState := d_user.UserState[testObs]{
		User: d_user.User{
			Identity: d_user.UserIdentity{AuthLevel: d_user.AuthRead},
		},
	}

	result := guard(userState, d_user.AuthWrite)
	if result == nil {
		t.Fatal("DefaultGuard should return a redirect when user has insufficient level")
	}

	redirect, ok := result.(*d_action.RedirectResponse)
	if !ok {
		t.Fatalf("DefaultGuard should return *RedirectResponse, got %T", result)
	}
	if redirect.TargetRoute != "denied_route" {
		t.Errorf("RedirectResponse.TargetRoute = %v, want denied_route", redirect.TargetRoute)
	}
}

func TestDefaultGuard_BlockedUser(t *testing.T) {
	guard := DefaultGuard[testObs]("blocked_page")

	userState := d_user.UserState[testObs]{
		User: d_user.User{
			Identity: d_user.UserIdentity{AuthLevel: d_user.AuthBlocked},
		},
	}

	result := guard(userState, d_user.AuthRead)
	if result == nil {
		t.Fatal("DefaultGuard should deny blocked user")
	}

	redirect, ok := result.(*d_action.RedirectResponse)
	if !ok {
		t.Fatalf("expected *RedirectResponse, got %T", result)
	}
	if redirect.TargetRoute != "blocked_page" {
		t.Errorf("TargetRoute = %v, want blocked_page", redirect.TargetRoute)
	}
}

func TestInternalGuard_AccessGranted(t *testing.T) {
	guard := InternalGuard[testObs]("denied_route")

	userState := d_user.UserState[testObs]{
		User: d_user.User{
			Internal: &d_user.UserInternal{
				Matricula: "12345",
				Cargo:     "Developer",
			},
		},
	}

	result := guard(userState, d_user.AuthRead)
	if result != nil {
		t.Errorf("InternalGuard should return nil when user has internal data, got %v", result)
	}
}

func TestInternalGuard_AccessDenied(t *testing.T) {
	guard := InternalGuard[testObs]("denied_route")

	userState := d_user.UserState[testObs]{
		User: d_user.User{
			Internal: nil,
		},
	}

	result := guard(userState, d_user.AuthRead)
	if result == nil {
		t.Fatal("InternalGuard should return a redirect when user has no internal data")
	}

	redirect, ok := result.(*d_action.RedirectResponse)
	if !ok {
		t.Fatalf("InternalGuard should return *RedirectResponse, got %T", result)
	}
	if redirect.TargetRoute != "denied_route" {
		t.Errorf("RedirectResponse.TargetRoute = %v, want denied_route", redirect.TargetRoute)
	}
}
