package main

import (
	"testing"

	"github.com/irissonnlima/chatgraph-go/chat"
)

// TestHandleStart validates the handleStart function behavior.
func TestHandleStart(t *testing.T) {
	engine := chat.NewEngine[Obs]()
	engine.RegisterRoute("start", handleStart)
	engine.RegisterRoute("menu", handleMenu)

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		User: chat.User{Name: "John"},
		Route: chat.Route{
			History:   []string{"start"},
			Separator: '/',
		},
		Observation: Obs{Field1: "initial", Field2: 0},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "hello"},
	}

	expectedActions := []chat.ExpectedAction{
		{Type: chat.ExecSendMessage, Message: &chat.Message{TextMessage: chat.TextMessage{Detail: "Hello John!"}}},
		{Type: chat.ExecSendMessage, Message: &chat.Message{TextMessage: chat.TextMessage{Detail: "Field1: initial, Field2: 0"}}},
		{Type: chat.ExecSetObservation, Observation: `{"field1":"updated","field2":1}`},
	}

	expectedReturn := chat.Route{
		History:   []string{"start", "menu"},
		Separator: '/',
	}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}

// TestHandleMenu_End validates handleMenu when user types "end".
func TestHandleMenu_End(t *testing.T) {
	engine := chat.NewEngine[Obs]()
	engine.RegisterRoute("menu", handleMenu)

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start", "menu"},
			Separator: '/',
		},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "end"},
	}

	expectedActions := []chat.ExpectedAction{}

	expectedReturn := chat.EndAction{ID: "session_ended"}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}

// TestHandleMenu_Start validates handleMenu when user types "start".
func TestHandleMenu_Start(t *testing.T) {
	engine := chat.NewEngine[Obs]()
	engine.RegisterRoute("menu", handleMenu)

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start", "menu"},
			Separator: '/',
		},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "start"},
	}

	expectedActions := []chat.ExpectedAction{}

	expectedReturn := chat.RedirectResponse{TargetRoute: "start"}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}

// TestHandleMenu_Default validates handleMenu when user types something else.
func TestHandleMenu_Default(t *testing.T) {
	engine := chat.NewEngine[Obs]()
	engine.RegisterRoute("menu", handleMenu)

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start", "menu"},
			Separator: '/',
		},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "something"},
	}

	expectedActions := []chat.ExpectedAction{
		{Type: chat.ExecSendMessage, Message: &chat.Message{TextMessage: chat.TextMessage{Detail: "Type 'end' to finish or 'start' to restart."}}},
	}

	// When handler returns nil, Engine returns route.Next(currentRoute) - staying on the same route
	expectedReturn := chat.Route{
		History:   []string{"start", "menu", "menu"},
		Separator: '/',
	}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}
