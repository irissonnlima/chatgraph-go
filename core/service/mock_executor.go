package service

import (
	"fmt"

	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_file "github.com/irissonnlima/chatgraph-go/core/domain/file"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

var ErrPrematureEnded = fmt.Errorf("session ended prematurely during testing")
var ErrPrematureTransfer = fmt.Errorf("session transferred prematurely during testing")

// mockExecutor is a mock executor that records actions.
type mockExecutor struct {
	expectedExec []ExpectedAction
}

func newMockExecutor() *mockExecutor {
	expectedExec := make([]ExpectedAction, 0)
	return &mockExecutor{
		expectedExec: expectedExec,
	}
}

func (m *mockExecutor) SendMessage(chatID d_user.ChatID, msg d_message.Message, platform string) error {
	m.expectedExec = append(m.expectedExec, ExpectedAction{
		Type:    ExecSendMessage,
		Message: &msg,
	})
	return nil
}

func (m *mockExecutor) SetObservation(chatID d_user.ChatID, observation string) error {
	m.expectedExec = append(m.expectedExec, ExpectedAction{
		Type:        ExecSetObservation,
		Observation: observation,
	})
	return nil
}

func (m *mockExecutor) EndSession(chatID d_user.ChatID, actionId string) error {
	return ErrPrematureEnded
}

func (m *mockExecutor) SetRoute(chatID d_user.ChatID, route string) error {
	m.expectedExec = append(m.expectedExec, ExpectedAction{
		Type:  ExecSetRoute,
		Route: route,
	})
	return nil
}

func (m *mockExecutor) TransferToMenu(chatID d_user.ChatID, transfer d_action.TransferToMenu, msg d_message.Message) error {
	return ErrPrematureTransfer
}

func (m *mockExecutor) UploadFile(filepath string) (*d_file.File, error) {
	m.expectedExec = append(m.expectedExec, ExpectedAction{
		Type:     ExecUploadFile,
		FilePath: filepath,
	})
	return &d_file.File{ID: "test-id", URL: "test-url", Name: filepath}, nil
}

func (m *mockExecutor) GetFile(fileID string) (*d_file.File, error) {
	m.expectedExec = append(m.expectedExec, ExpectedAction{
		Type:   ExecGetFile,
		FileID: fileID,
	})
	return &d_file.File{ID: fileID, URL: "test-url", Name: "test"}, nil
}
