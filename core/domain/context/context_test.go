package d_context

import (
	"errors"
	"testing"
	"time"

	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_file "github.com/irissonnlima/chatgraph-go/core/domain/file"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_route "github.com/irissonnlima/chatgraph-go/core/domain/route"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

// MockRouter implements IBotExecutor for testing
type MockRouter struct {
	SendMessageFunc    func(chatID d_user.ChatID, message d_message.Message, platform string) error
	SetObservationFunc func(chatID d_user.ChatID, observation string) error
	EndSessionFunc     func(chatID d_user.ChatID, actionId string) error
	SetRouteFunc       func(chatID d_user.ChatID, route string) error
	TransferFunc       func(chatID d_user.ChatID, transfer d_action.TransferToMenu, message d_message.Message) error
	UploadFileFunc     func(filepath string) (*d_file.File, error)
	GetFileFunc        func(fileID string) (*d_file.File, error)
}

func (m *MockRouter) SendMessage(chatID d_user.ChatID, message d_message.Message, platform string) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(chatID, message, platform)
	}
	return nil
}

func (m *MockRouter) SetObservation(chatID d_user.ChatID, observation string) error {
	if m.SetObservationFunc != nil {
		return m.SetObservationFunc(chatID, observation)
	}
	return nil
}

func (m *MockRouter) EndSession(chatID d_user.ChatID, actionId string) error {
	if m.EndSessionFunc != nil {
		return m.EndSessionFunc(chatID, actionId)
	}
	return nil
}

func (m *MockRouter) SetRoute(chatID d_user.ChatID, route string) error {
	if m.SetRouteFunc != nil {
		return m.SetRouteFunc(chatID, route)
	}
	return nil
}

func (m *MockRouter) TransferToMenu(chatID d_user.ChatID, transfer d_action.TransferToMenu, message d_message.Message) error {
	if m.TransferFunc != nil {
		return m.TransferFunc(chatID, transfer, message)
	}
	return nil
}

func (m *MockRouter) UploadFile(filepath string) (*d_file.File, error) {
	if m.UploadFileFunc != nil {
		return m.UploadFileFunc(filepath)
	}
	return nil, nil
}

func (m *MockRouter) GetFile(fileID string) (*d_file.File, error) {
	if m.GetFileFunc != nil {
		return m.GetFileFunc(fileID)
	}
	return nil, nil
}

type TestObservation struct {
	Value string `json:"value"`
}

func TestNewChatContext(t *testing.T) {
	userState := d_user.UserState[TestObservation]{
		SessionID: 123,
		ChatID:    d_user.ChatID{UserID: "user1", CompanyID: "comp1"},
	}
	message := d_message.Message{
		TextMessage: d_message.TextMessage{Detail: "Hello"},
	}
	router := &MockRouter{}
	timeout := 5 * time.Second

	ctx, cancel := NewChatContext(userState, message, router, timeout)
	defer cancel()

	if ctx.UserState.SessionID != 123 {
		t.Errorf("NewChatContext().UserState.SessionID = %v, want %v", ctx.UserState.SessionID, 123)
	}
	if ctx.Message.TextMessage.Detail != "Hello" {
		t.Errorf("NewChatContext().Message.TextMessage.Detail = %v, want %v", ctx.Message.TextMessage.Detail, "Hello")
	}
	if ctx.Context == nil {
		t.Error("NewChatContext().Context should not be nil")
	}
}

func TestChatContext_SendMessage(t *testing.T) {
	var sentMessage d_message.Message
	var sentPlatform string

	router := &MockRouter{
		SendMessageFunc: func(chatID d_user.ChatID, message d_message.Message, platform string) error {
			sentMessage = message
			sentPlatform = platform
			return nil
		},
	}

	userState := d_user.UserState[TestObservation]{
		ChatID:   d_user.ChatID{UserID: "user1"},
		Platform: "whatsapp",
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	defer cancel()

	msg := d_message.Message{TextMessage: d_message.TextMessage{Detail: "Test"}}
	err := ctx.SendMessage(msg)

	if err != nil {
		t.Errorf("SendMessage() error = %v", err)
	}
	if sentMessage.TextMessage.Detail != "Test" {
		t.Errorf("SendMessage() sent detail = %v, want %v", sentMessage.TextMessage.Detail, "Test")
	}
	if sentPlatform != "whatsapp" {
		t.Errorf("SendMessage() sent platform = %v, want %v", sentPlatform, "whatsapp")
	}
}

func TestChatContext_SendMessage_ContextCanceled(t *testing.T) {
	router := &MockRouter{}
	userState := d_user.UserState[TestObservation]{}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	cancel() // Cancel immediately

	err := ctx.SendMessage(d_message.Message{})
	if err == nil {
		t.Error("SendMessage() should return error when context is canceled")
	}
}

func TestChatContext_SendTextMessage(t *testing.T) {
	var sentMessage d_message.Message

	router := &MockRouter{
		SendMessageFunc: func(chatID d_user.ChatID, message d_message.Message, platform string) error {
			sentMessage = message
			return nil
		},
	}

	userState := d_user.UserState[TestObservation]{
		ChatID: d_user.ChatID{UserID: "user1"},
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	defer cancel()

	err := ctx.SendTextMessage("Hello World")

	if err != nil {
		t.Errorf("SendTextMessage() error = %v", err)
	}
	if sentMessage.TextMessage.Detail != "Hello World" {
		t.Errorf("SendTextMessage() sent detail = %v, want %v", sentMessage.TextMessage.Detail, "Hello World")
	}
}

func TestChatContext_GetObservation(t *testing.T) {
	obs := TestObservation{Value: "test_value"}
	userState := d_user.UserState[TestObservation]{
		Observation: obs,
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, &MockRouter{}, 5*time.Second)
	defer cancel()

	got := ctx.GetObservation()
	if got.Value != "test_value" {
		t.Errorf("GetObservation().Value = %v, want %v", got.Value, "test_value")
	}
}

func TestChatContext_SetObservation(t *testing.T) {
	var savedObs string

	router := &MockRouter{
		SetObservationFunc: func(chatID d_user.ChatID, observation string) error {
			savedObs = observation
			return nil
		},
	}

	userState := d_user.UserState[TestObservation]{
		ChatID: d_user.ChatID{UserID: "user1"},
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	defer cancel()

	err := ctx.SetObservation(TestObservation{Value: "new_value"})

	if err != nil {
		t.Errorf("SetObservation() error = %v", err)
	}
	if savedObs != `{"value":"new_value"}` {
		t.Errorf("SetObservation() saved = %v, want %v", savedObs, `{"value":"new_value"}`)
	}
}

func TestChatContext_SetObservation_ContextCanceled(t *testing.T) {
	router := &MockRouter{}
	userState := d_user.UserState[TestObservation]{}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	cancel()

	err := ctx.SetObservation(TestObservation{Value: "test"})
	if err == nil {
		t.Error("SetObservation() should return error when context is canceled")
	}
}

func TestChatContext_SetObservation_RouterError(t *testing.T) {
	router := &MockRouter{
		SetObservationFunc: func(chatID d_user.ChatID, observation string) error {
			return errors.New("router error")
		},
	}

	userState := d_user.UserState[TestObservation]{
		ChatID: d_user.ChatID{UserID: "user1"},
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	defer cancel()

	err := ctx.SetObservation(TestObservation{Value: "test"})
	if err == nil {
		t.Error("SetObservation() should return error when router fails")
	}
}

func TestChatContext_GetRoute(t *testing.T) {
	route := d_route.NewRoute("menu.submenu.action", '.')
	userState := d_user.UserState[TestObservation]{
		Route: route,
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, &MockRouter{}, 5*time.Second)
	defer cancel()

	got := ctx.GetRoute()
	if got.Current() != route.Current() {
		t.Errorf("GetRoute().Current() = %v, want %v", got.Current(), route.Current())
	}
}

func TestChatContext_NextRoute(t *testing.T) {
	route := d_route.NewRoute("menu.submenu", '.')
	userState := d_user.UserState[TestObservation]{
		Route: route,
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, &MockRouter{}, 5*time.Second)
	defer cancel()

	next := ctx.NextRoute("action")
	if next.Current() != "action" {
		t.Errorf("NextRoute().Current() = %v, want %v", next.Current(), "action")
	}
}

func TestChatContext_NextRoute_ContextCanceled(t *testing.T) {
	route := d_route.NewRoute("menu", '.')
	userState := d_user.UserState[TestObservation]{
		Route: route,
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, &MockRouter{}, 5*time.Second)
	cancel()

	next := ctx.NextRoute("action")
	// Should return the original route when context is canceled
	if next.Current() != route.Current() {
		t.Errorf("NextRoute() when canceled Current() = %v, want %v", next.Current(), route.Current())
	}
}

func TestChatContext_LoadFile(t *testing.T) {
	expectedFile := &d_file.File{
		ID:   "f1",
		Name: "test.txt",
		URL:  "https://example.com/test.txt",
	}

	router := &MockRouter{
		UploadFileFunc: func(filepath string) (*d_file.File, error) {
			return expectedFile, nil
		},
	}

	userState := d_user.UserState[TestObservation]{
		ChatID: d_user.ChatID{UserID: "user1"},
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	defer cancel()

	got, err := ctx.LoadFile("/path/to/file.txt")
	if err != nil {
		t.Errorf("LoadFile() error = %v", err)
	}
	if got.ID != expectedFile.ID {
		t.Errorf("LoadFile().ID = %v, want %v", got.ID, expectedFile.ID)
	}
}

func TestChatContext_LoadFile_ContextCanceled(t *testing.T) {
	router := &MockRouter{}
	userState := d_user.UserState[TestObservation]{}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	cancel()

	_, err := ctx.LoadFile("/path/to/file.txt")
	if err == nil {
		t.Error("LoadFile() should return error when context is canceled")
	}
}

func TestChatContext_LoadFileBytes(t *testing.T) {
	expectedFile := &d_file.File{
		ID:   "f1",
		Name: "test.pdf",
		URL:  "https://example.com/test.pdf",
	}

	router := &MockRouter{
		UploadFileFunc: func(filepath string) (*d_file.File, error) {
			return expectedFile, nil
		},
	}

	userState := d_user.UserState[TestObservation]{
		ChatID: d_user.ChatID{UserID: "user1"},
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	defer cancel()

	data := []byte("test content")
	got, err := ctx.LoadFileBytes("report.pdf", data)
	if err != nil {
		t.Errorf("LoadFileBytes() error = %v", err)
	}
	if got.ID != expectedFile.ID {
		t.Errorf("LoadFileBytes().ID = %v, want %v", got.ID, expectedFile.ID)
	}
}

func TestChatContext_LoadFileBytes_ContextCanceled(t *testing.T) {
	router := &MockRouter{}
	userState := d_user.UserState[TestObservation]{}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	cancel()

	_, err := ctx.LoadFileBytes("test.txt", []byte("data"))
	if err == nil {
		t.Error("LoadFileBytes() should return error when context is canceled")
	}
}

func TestChatContext_LoadFileBytes_UploadError(t *testing.T) {
	router := &MockRouter{
		UploadFileFunc: func(filepath string) (*d_file.File, error) {
			return nil, errors.New("upload failed")
		},
	}

	userState := d_user.UserState[TestObservation]{
		ChatID: d_user.ChatID{UserID: "user1"},
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	defer cancel()

	_, err := ctx.LoadFileBytes("test.txt", []byte("data"))
	if err == nil {
		t.Error("LoadFileBytes() should return error when upload fails")
	}
}

func TestChatContext_GetFile(t *testing.T) {
	expectedFile := &d_file.File{
		ID:   "f1",
		Name: "test.txt",
		URL:  "https://example.com/test.txt",
	}

	router := &MockRouter{
		GetFileFunc: func(fileID string) (*d_file.File, error) {
			if fileID == "f1" {
				return expectedFile, nil
			}
			return nil, errors.New("file not found")
		},
	}

	userState := d_user.UserState[TestObservation]{
		ChatID: d_user.ChatID{UserID: "user1"},
	}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	defer cancel()

	got, err := ctx.GetFile("f1")
	if err != nil {
		t.Errorf("GetFile() error = %v", err)
	}
	if got.ID != expectedFile.ID {
		t.Errorf("GetFile().ID = %v, want %v", got.ID, expectedFile.ID)
	}
}

func TestChatContext_GetFile_ContextCanceled(t *testing.T) {
	router := &MockRouter{}
	userState := d_user.UserState[TestObservation]{}

	ctx, cancel := NewChatContext(userState, d_message.Message{}, router, 5*time.Second)
	cancel()

	_, err := ctx.GetFile("f1")
	if err == nil {
		t.Error("GetFile() should return error when context is canceled")
	}
}
