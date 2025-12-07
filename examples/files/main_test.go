package main

import (
	"testing"

	"github.com/irissonnlima/chatgraph-go/chat"
)

// TestStart validates the start route sends instructions.
func TestStart(t *testing.T) {
	engine := chat.NewEngine[Obs]()

	engine.RegisterRoute("start", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendTextMessage("Send 'file' to receive a file, or 'upload' to upload from bytes.")
		return ctx.NextRoute("handle_input")
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
		{Type: chat.ExecSendMessage, Message: &chat.Message{TextMessage: chat.TextMessage{Detail: "Send 'file' to receive a file, or 'upload' to upload from bytes."}}},
	}

	expectedReturn := chat.Route{
		History:   []string{"start", "handle_input"},
		Separator: '/',
	}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}

// TestHandleInput_File validates the "file" command triggers GetFile action.
func TestHandleInput_File(t *testing.T) {
	engine := chat.NewEngine[Obs]()

	engine.RegisterRoute("handle_input", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		input := ctx.Message.EntireText()

		switch input {
		case "file":
			file, err := ctx.LoadFile("README.md")
			if err != nil {
				ctx.SendTextMessage("Error loading file: " + err.Error())
				return nil
			}
			if file != nil {
				ctx.SendMessage(chat.Message{
					TextMessage: chat.TextMessage{
						Detail: "Here's your file:",
					},
					File: *file,
				})
			}
		case "upload":
			content := []byte("Hello, this is a test file content!")
			file, err := ctx.LoadFileBytes("test-file.txt", content)
			if err != nil {
				ctx.SendTextMessage("Error uploading file: " + err.Error())
				return nil
			}
			if file != nil {
				ctx.SendMessage(chat.Message{
					TextMessage: chat.TextMessage{
						Detail: "File created from bytes:",
					},
					File: *file,
				})
			}
		default:
			ctx.SendTextMessage("Unknown command. Try 'file' or 'upload'.")
		}

		return nil
	})

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start", "handle_input"},
			Separator: '/',
		},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "file"},
	}

	// The mock will return a file with ID "test-id"
	expectedActions := []chat.ExpectedAction{
		{Type: chat.ExecUploadFile, FilePath: "README.md"},
		{Type: chat.ExecSendMessage, Message: &chat.Message{
			TextMessage: chat.TextMessage{Detail: "Here's your file:"},
			File:        chat.File{ID: "test-id", URL: "test-url", Name: "README.md"},
		}},
	}

	// When handler returns nil, Engine returns route.Next(currentRoute) - staying on the same route
	expectedReturn := chat.Route{
		History:   []string{"start", "handle_input", "handle_input"},
		Separator: '/',
	}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}

// TestHandleInput_Upload validates the "upload" command triggers UploadFile action.
func TestHandleInput_Upload(t *testing.T) {
	engine := chat.NewEngine[Obs]()

	engine.RegisterRoute("handle_input", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		input := ctx.Message.EntireText()

		switch input {
		case "file":
			file, err := ctx.LoadFile("README.md")
			if err != nil {
				ctx.SendTextMessage("Error loading file: " + err.Error())
				return nil
			}
			if file != nil {
				ctx.SendMessage(chat.Message{
					TextMessage: chat.TextMessage{
						Detail: "Here's your file:",
					},
					File: *file,
				})
			}
		case "upload":
			content := []byte("Hello, this is a test file content!")
			file, err := ctx.LoadFileBytes("test-file.txt", content)
			if err != nil {
				ctx.SendTextMessage("Error uploading file: " + err.Error())
				return nil
			}
			if file != nil {
				ctx.SendMessage(chat.Message{
					TextMessage: chat.TextMessage{
						Detail: "File created from bytes:",
					},
					File: *file,
				})
			}
		default:
			ctx.SendTextMessage("Unknown command. Try 'file' or 'upload'.")
		}

		return nil
	})

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start", "handle_input"},
			Separator: '/',
		},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "upload"},
	}

	// The mock will return a file with ID "test-id"
	// Note: LoadFileBytes creates a temp file with random name, so we can't match exact File.Name
	// We only validate the action types here, not the exact message content
	expectedActions := []chat.ExpectedAction{
		{Type: chat.ExecUploadFile},  // FilePath is temp file, can't predict exact name
		{Type: chat.ExecSendMessage}, // File.Name includes temp path, can't match exactly
	}

	// When handler returns nil, Engine returns route.Next(currentRoute) - staying on the same route
	expectedReturn := chat.Route{
		History:   []string{"start", "handle_input", "handle_input"},
		Separator: '/',
	}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}

// TestHandleInput_Unknown validates unknown commands.
func TestHandleInput_Unknown(t *testing.T) {
	engine := chat.NewEngine[Obs]()

	engine.RegisterRoute("handle_input", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		input := ctx.Message.EntireText()

		switch input {
		case "file":
			file, err := ctx.LoadFile("README.md")
			if err != nil {
				ctx.SendTextMessage("Error loading file: " + err.Error())
				return nil
			}
			if file != nil {
				ctx.SendMessage(chat.Message{
					TextMessage: chat.TextMessage{
						Detail: "Here's your file:",
					},
					File: *file,
				})
			}
		case "upload":
			content := []byte("Hello, this is a test file content!")
			file, err := ctx.LoadFileBytes("test-file.txt", content)
			if err != nil {
				ctx.SendTextMessage("Error uploading file: " + err.Error())
				return nil
			}
			if file != nil {
				ctx.SendMessage(chat.Message{
					TextMessage: chat.TextMessage{
						Detail: "File created from bytes:",
					},
					File: *file,
				})
			}
		default:
			ctx.SendTextMessage("Unknown command. Try 'file' or 'upload'.")
		}

		return nil
	})

	tester := chat.NewEngineTester(t, engine)

	userState := chat.UserState[Obs]{
		Route: chat.Route{
			History:   []string{"start", "handle_input"},
			Separator: '/',
		},
	}

	msg := chat.Message{
		TextMessage: chat.TextMessage{Detail: "something"},
	}

	expectedActions := []chat.ExpectedAction{
		{Type: chat.ExecSendMessage, Message: &chat.Message{TextMessage: chat.TextMessage{Detail: "Unknown command. Try 'file' or 'upload'."}}},
	}

	// When handler returns nil, Engine returns route.Next(currentRoute) - staying on the same route
	expectedReturn := chat.Route{
		History:   []string{"start", "handle_input", "handle_input"},
		Separator: '/',
	}

	tester.Execute(userState, msg, expectedActions, expectedReturn)
}
