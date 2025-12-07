package main

import (
	"testing"

	"github.com/irissonnlima/chatgraph-go/chat"
)

// TestStart validates the start handler sends a message with buttons.
func TestStart(t *testing.T) {
	engine := chat.NewEngine[Obs]()

	engine.RegisterRoute("start", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendMessage(chat.Message{
			TextMessage: chat.TextMessage{
				Title:  "Welcome!",
				Detail: "Please choose an option:",
			},
			Buttons: []chat.Button{
				{
					Type:   chat.POSTBACK,
					Title:  "Option A",
					Detail: "option_a",
				},
				{
					Type:   chat.POSTBACK,
					Title:  "Option B",
					Detail: "option_b",
				},
				{
					Type:   chat.URL,
					Title:  "Visit Website",
					Detail: "https://example.com",
				},
			},
		})
		return ctx.NextRoute("handle_choice")
	})

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start"},
			Separator: '/',
		},
	}

	msg := chat.Message{}

	expectedActions := []chat.ExpectedAction{
		{
			Type: chat.ExecSendMessage,
			Message: &chat.Message{
				TextMessage: chat.TextMessage{
					Title:  "Welcome!",
					Detail: "Please choose an option:",
				},
				Buttons: []chat.Button{
					{Type: chat.POSTBACK, Title: "Option A", Detail: "option_a"},
					{Type: chat.POSTBACK, Title: "Option B", Detail: "option_b"},
					{Type: chat.URL, Title: "Visit Website", Detail: "https://example.com"},
				},
			},
		},
	}

	expectedReturn := chat.Route{
		History:   []string{"start", "handle_choice"},
		Separator: '/',
	}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}

// TestHandleChoice_OptionA validates selecting option A.
func TestHandleChoice_OptionA(t *testing.T) {
	engine := chat.NewEngine[Obs]()

	engine.RegisterRoute("handle_choice", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		choice := ctx.Message.EntireText()

		switch choice {
		case "option_a":
			ctx.SendTextMessage("You selected Option A!")
		case "option_b":
			ctx.SendTextMessage("You selected Option B!")
		default:
			ctx.SendTextMessage("Unknown option: " + choice)
		}

		return &chat.RedirectResponse{TargetRoute: "start"}
	})

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start", "handle_choice"},
			Separator: '/',
		},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "option_a"},
	}

	expectedActions := []chat.ExpectedAction{
		{Type: chat.ExecSendMessage, Message: &chat.Message{TextMessage: chat.TextMessage{Detail: "You selected Option A!"}}},
	}

	expectedReturn := &chat.RedirectResponse{TargetRoute: "start"}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}

// TestHandleChoice_OptionB validates selecting option B.
func TestHandleChoice_OptionB(t *testing.T) {
	engine := chat.NewEngine[Obs]()

	engine.RegisterRoute("handle_choice", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		choice := ctx.Message.EntireText()

		switch choice {
		case "option_a":
			ctx.SendTextMessage("You selected Option A!")
		case "option_b":
			ctx.SendTextMessage("You selected Option B!")
		default:
			ctx.SendTextMessage("Unknown option: " + choice)
		}

		return &chat.RedirectResponse{TargetRoute: "start"}
	})

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start", "handle_choice"},
			Separator: '/',
		},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "option_b"},
	}

	expectedActions := []chat.ExpectedAction{
		{Type: chat.ExecSendMessage, Message: &chat.Message{TextMessage: chat.TextMessage{Detail: "You selected Option B!"}}},
	}

	expectedReturn := &chat.RedirectResponse{TargetRoute: "start"}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}

// TestHandleChoice_Unknown validates unknown input.
func TestHandleChoice_Unknown(t *testing.T) {
	engine := chat.NewEngine[Obs]()

	engine.RegisterRoute("handle_choice", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		choice := ctx.Message.EntireText()

		switch choice {
		case "option_a":
			ctx.SendTextMessage("You selected Option A!")
		case "option_b":
			ctx.SendTextMessage("You selected Option B!")
		default:
			ctx.SendTextMessage("Unknown option: " + choice)
		}

		return &chat.RedirectResponse{TargetRoute: "start"}
	})

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start", "handle_choice"},
			Separator: '/',
		},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "invalid"},
	}

	expectedActions := []chat.ExpectedAction{
		{Type: chat.ExecSendMessage, Message: &chat.Message{TextMessage: chat.TextMessage{Detail: "Unknown option: invalid"}}},
	}

	expectedReturn := &chat.RedirectResponse{TargetRoute: "start"}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}
