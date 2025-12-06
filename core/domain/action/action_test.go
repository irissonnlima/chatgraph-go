package d_action

import (
	"testing"
)

func TestEndAction_IsRouteReturn(t *testing.T) {
	action := EndAction{ID: "test"}
	// Should not panic - just verify it implements the interface
	action.IsRouteReturn()
}

func TestRedirectResponse_IsRouteReturn(t *testing.T) {
	redirect := RedirectResponse{TargetRoute: "start"}
	// Should not panic - just verify it implements the interface
	redirect.IsRouteReturn()
}

func TestTransferToMenu_IsRouteReturn(t *testing.T) {
	transfer := TransferToMenu{MenuID: 1, Route: "menu"}
	// Should not panic - just verify it implements the interface
	transfer.IsRouteReturn()
}

func TestEndAction_Fields(t *testing.T) {
	action := EndAction{
		ID:           "end_123",
		Name:         "Session End",
		DepartmentID: 5,
		Observation:  "User requested end",
		LastUpdate:   "2024-01-01",
	}

	if action.ID != "end_123" {
		t.Errorf("EndAction.ID = %v, want end_123", action.ID)
	}
	if action.Name != "Session End" {
		t.Errorf("EndAction.Name = %v, want Session End", action.Name)
	}
	if action.DepartmentID != 5 {
		t.Errorf("EndAction.DepartmentID = %v, want 5", action.DepartmentID)
	}
	if action.Observation != "User requested end" {
		t.Errorf("EndAction.Observation = %v, want User requested end", action.Observation)
	}
	if action.LastUpdate != "2024-01-01" {
		t.Errorf("EndAction.LastUpdate = %v, want 2024-01-01", action.LastUpdate)
	}
}

func TestRedirectResponse_Fields(t *testing.T) {
	redirect := RedirectResponse{TargetRoute: "menu"}

	if redirect.TargetRoute != "menu" {
		t.Errorf("RedirectResponse.TargetRoute = %v, want menu", redirect.TargetRoute)
	}
}

func TestTransferToMenu_Fields(t *testing.T) {
	transfer := TransferToMenu{MenuID: 10, Route: "support"}

	if transfer.MenuID != 10 {
		t.Errorf("TransferToMenu.MenuID = %v, want 10", transfer.MenuID)
	}
	if transfer.Route != "support" {
		t.Errorf("TransferToMenu.Route = %v, want support", transfer.Route)
	}
}
