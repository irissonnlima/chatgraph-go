package service

import (
	"testing"
	"time"

	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_file "github.com/irissonnlima/chatgraph-go/core/domain/file"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_router "github.com/irissonnlima/chatgraph-go/core/domain/router"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

type TestObs struct {
	Value string
}

// MockMessageReceiver for testing
type MockMessageReceiver struct{}

func (m *MockMessageReceiver) ConsumeMessage() <-chan struct {
	UserState d_user.UserState[TestObs]
	Message   d_message.Message
} {
	ch := make(chan struct {
		UserState d_user.UserState[TestObs]
		Message   d_message.Message
	})
	return ch
}

// MockBotExecutor for testing
type MockBotExecutor struct{}

func (m *MockBotExecutor) SendMessage(to d_user.ChatID, message d_message.Message, platform string) error {
	return nil
}

func (m *MockBotExecutor) SetObservation(chatID d_user.ChatID, observation string) error {
	return nil
}

func (m *MockBotExecutor) SetRoute(chatID d_user.ChatID, route string) error {
	return nil
}

func (m *MockBotExecutor) EndSession(chatID d_user.ChatID, actionId string) error {
	return nil
}

func (m *MockBotExecutor) TransferToMenu(chatID d_user.ChatID, transfer d_action.TransferToMenu, message d_message.Message) error {
	return nil
}

func (m *MockBotExecutor) UploadFile(filepath string) (*d_file.File, error) {
	return nil, nil
}

func (m *MockBotExecutor) GetFile(fileID string) (*d_file.File, error) {
	return nil, nil
}

func TestNewChatbotApp(t *testing.T) {
	receiver := &MockMessageReceiver{}
	executor := &MockBotExecutor{}

	app := NewChatbotApp[TestObs](receiver, executor)

	if app == nil {
		t.Error("NewChatbotApp() should return non-nil app")
	}
	if app.routes == nil {
		t.Error("NewChatbotApp().routes should be initialized")
	}
}

func TestNewChatbotApp_WithDefaultOptions(t *testing.T) {
	receiver := &MockMessageReceiver{}
	executor := &MockBotExecutor{}

	timeout := d_router.TimeoutRouteOps{
		Duration: 10 * time.Minute,
		Route:    "custom_timeout",
	}
	loopCount := d_router.LoopCountRouteOps{
		Count: 5,
		Route: "custom_loop",
	}

	opts := d_router.RouterHandlerOptions{
		Timeout:   &timeout,
		LoopCount: &loopCount,
	}

	app := NewChatbotApp[TestObs](receiver, executor, opts)

	if app == nil {
		t.Error("NewChatbotApp() should return non-nil app")
	}
	if app.defaultOptions.Timeout.Duration != 10*time.Minute {
		t.Errorf("NewChatbotApp().defaultOptions.Timeout.Duration = %v, want 10m", app.defaultOptions.Timeout.Duration)
	}
	if app.defaultOptions.LoopCount.Count != 5 {
		t.Errorf("NewChatbotApp().defaultOptions.LoopCount.Count = %v, want 5", app.defaultOptions.LoopCount.Count)
	}
}

func TestChatbotApp_applyTriggers(t *testing.T) {
	receiver := &MockMessageReceiver{}
	executor := &MockBotExecutor{}

	app := NewChatbotApp[TestObs](receiver, executor)
	app.routeTriggers = []d_router.RouteTrigger{
		{Regex: "^menu$", Route: "menu_route"},
		{Regex: "^help", Route: "help_route"},
		{Regex: "\\d+", Route: "number_route"},
	}

	tests := []struct {
		name        string
		messageText string
		want        string
	}{
		{"exact match", "menu", "menu_route"},
		{"prefix match", "help me", "help_route"},
		{"pattern match", "call 123", "number_route"},
		{"no match", "random text", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.applyTriggers(tt.messageText)
			if got != tt.want {
				t.Errorf("applyTriggers(%q) = %v, want %v", tt.messageText, got, tt.want)
			}
		})
	}
}

func TestChatbotApp_applyTriggers_InvalidRegex(t *testing.T) {
	receiver := &MockMessageReceiver{}
	executor := &MockBotExecutor{}

	app := NewChatbotApp[TestObs](receiver, executor)
	app.routeTriggers = []d_router.RouteTrigger{
		{Regex: "[invalid", Route: "invalid_route"}, // invalid regex
		{Regex: "^valid$", Route: "valid_route"},
	}

	// Should skip invalid regex and still check others
	got := app.applyTriggers("valid")
	if got != "valid_route" {
		t.Errorf("applyTriggers() = %v, want valid_route", got)
	}
}
