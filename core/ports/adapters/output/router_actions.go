// Package adapter_output defines service interfaces for the chatbot's core functionality.
// These interfaces abstract the routing, messaging, and session management operations.
package adapter_output

import (
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_file "github.com/irissonnlima/chatgraph-go/core/domain/file"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

// RouterService defines the interface for routing and messaging operations.
// Implementations handle message sending, session management, and state updates.
type RouterService interface {
	// SendMessage sends a message to the specified chat ID.
	// Returns an error if the message could not be delivered.
	SendMessage(to d_user.ChatID, message d_message.Message, platform string) error

	// SetObservation updates the observation data for the specified chat.
	// The observation is stored as a JSON string.
	SetObservation(chatID d_user.ChatID, observation string) error

	// SetRoute updates the current route for the specified chat.
	// This determines which handler will process the next message.
	SetRoute(chatID d_user.ChatID, route string) error

	// EndSession terminates the session for the specified chat.
	// This should clean up any session-related resources.
	EndSession(chatID d_user.ChatID, actionId string) error

	// TransferToMenu transfers the user to other menu.
	TransferToMenu(chatID d_user.ChatID, transfer d_action.TransferToMenu, message d_message.Message) error

	// UploadFile uploads a file from the given filepath.
	UploadFile(filepath string) (*d_file.File, error)

	// GetFile retrieves a file by its ID.
	GetFile(fileID string) (*d_file.File, error)
}
