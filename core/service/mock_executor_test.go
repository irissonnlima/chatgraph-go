package service

import (
	"testing"

	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

// TestMockExecutor_SendMessage tests SendMessage recording.
func TestMockExecutor_SendMessage(t *testing.T) {
	mock := newMockExecutor()

	msg := d_message.Message{
		TextMessage: d_message.TextMessage{Detail: "Hello"},
	}

	err := mock.SendMessage(d_user.ChatID{}, msg, "whatsapp")
	if err != nil {
		t.Fatalf("SendMessage returned error: %v", err)
	}

	if len(mock.expectedExec) != 1 {
		t.Fatalf("expected 1 action, got %d", len(mock.expectedExec))
	}

	if mock.expectedExec[0].Type != ExecSendMessage {
		t.Errorf("expected ExecSendMessage, got %d", mock.expectedExec[0].Type)
	}
}

// TestMockExecutor_SetObservation tests SetObservation recording.
func TestMockExecutor_SetObservation(t *testing.T) {
	mock := newMockExecutor()

	err := mock.SetObservation(d_user.ChatID{}, `{"key":"value"}`)
	if err != nil {
		t.Fatalf("SetObservation returned error: %v", err)
	}

	if len(mock.expectedExec) != 1 {
		t.Fatalf("expected 1 action, got %d", len(mock.expectedExec))
	}

	if mock.expectedExec[0].Type != ExecSetObservation {
		t.Errorf("expected ExecSetObservation, got %d", mock.expectedExec[0].Type)
	}

	if mock.expectedExec[0].Observation != `{"key":"value"}` {
		t.Errorf("expected observation '{\"key\":\"value\"}', got %s", mock.expectedExec[0].Observation)
	}
}

// TestMockExecutor_SetRoute tests SetRoute recording.
func TestMockExecutor_SetRoute(t *testing.T) {
	mock := newMockExecutor()

	err := mock.SetRoute(d_user.ChatID{}, "next_route")
	if err != nil {
		t.Fatalf("SetRoute returned error: %v", err)
	}

	if len(mock.expectedExec) != 1 {
		t.Fatalf("expected 1 action, got %d", len(mock.expectedExec))
	}

	if mock.expectedExec[0].Type != ExecSetRoute {
		t.Errorf("expected ExecSetRoute, got %d", mock.expectedExec[0].Type)
	}

	if mock.expectedExec[0].Route != "next_route" {
		t.Errorf("expected route 'next_route', got %s", mock.expectedExec[0].Route)
	}
}

// TestMockExecutor_UploadFile tests UploadFile recording.
func TestMockExecutor_UploadFile(t *testing.T) {
	mock := newMockExecutor()

	file, err := mock.UploadFile("/path/to/file.txt")
	if err != nil {
		t.Fatalf("UploadFile returned error: %v", err)
	}

	if file == nil {
		t.Fatal("UploadFile returned nil file")
	}

	if file.ID != "test-id" {
		t.Errorf("expected file ID 'test-id', got %s", file.ID)
	}

	if len(mock.expectedExec) != 1 {
		t.Fatalf("expected 1 action, got %d", len(mock.expectedExec))
	}

	if mock.expectedExec[0].Type != ExecUploadFile {
		t.Errorf("expected ExecUploadFile, got %d", mock.expectedExec[0].Type)
	}

	if mock.expectedExec[0].FilePath != "/path/to/file.txt" {
		t.Errorf("expected filePath '/path/to/file.txt', got %s", mock.expectedExec[0].FilePath)
	}
}

// TestMockExecutor_GetFile tests GetFile recording.
func TestMockExecutor_GetFile(t *testing.T) {
	mock := newMockExecutor()

	file, err := mock.GetFile("file-123")
	if err != nil {
		t.Fatalf("GetFile returned error: %v", err)
	}

	if file == nil {
		t.Fatal("GetFile returned nil file")
	}

	if file.ID != "file-123" {
		t.Errorf("expected file ID 'file-123', got %s", file.ID)
	}

	if len(mock.expectedExec) != 1 {
		t.Fatalf("expected 1 action, got %d", len(mock.expectedExec))
	}

	if mock.expectedExec[0].Type != ExecGetFile {
		t.Errorf("expected ExecGetFile, got %d", mock.expectedExec[0].Type)
	}

	if mock.expectedExec[0].FileID != "file-123" {
		t.Errorf("expected fileID 'file-123', got %s", mock.expectedExec[0].FileID)
	}
}

// TestMockExecutor_EndSession tests EndSession returns error.
func TestMockExecutor_EndSession(t *testing.T) {
	mock := newMockExecutor()

	err := mock.EndSession(d_user.ChatID{}, "session-ended")

	if err != ErrPrematureEnded {
		t.Errorf("expected ErrPrematureEnded, got %v", err)
	}
}

// TestMockExecutor_TransferToMenu tests TransferToMenu returns error.
func TestMockExecutor_TransferToMenu(t *testing.T) {
	mock := newMockExecutor()

	err := mock.TransferToMenu(d_user.ChatID{}, d_action.TransferToMenu{}, d_message.Message{})

	if err != ErrPrematureTransfer {
		t.Errorf("expected ErrPrematureTransfer, got %v", err)
	}
}

// TestMockExecutor_MultipleActions tests recording multiple actions.
func TestMockExecutor_MultipleActions(t *testing.T) {
	mock := newMockExecutor()

	mock.SendMessage(d_user.ChatID{}, d_message.Message{}, "whatsapp")
	mock.SetObservation(d_user.ChatID{}, "{}")
	mock.SetRoute(d_user.ChatID{}, "next")

	if len(mock.expectedExec) != 3 {
		t.Errorf("expected 3 actions, got %d", len(mock.expectedExec))
	}

	if mock.expectedExec[0].Type != ExecSendMessage {
		t.Errorf("action 0: expected ExecSendMessage")
	}
	if mock.expectedExec[1].Type != ExecSetObservation {
		t.Errorf("action 1: expected ExecSetObservation")
	}
	if mock.expectedExec[2].Type != ExecSetRoute {
		t.Errorf("action 2: expected ExecSetRoute")
	}
}
