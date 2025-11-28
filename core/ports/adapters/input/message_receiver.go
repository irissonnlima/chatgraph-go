// Package adapter_input defines input adapter interfaces for external services.
// These interfaces abstract the communication with external systems like message queues.
package adapter_input

import (
	d_message "chatgraph/core/domain/message"
	d_user "chatgraph/core/domain/user"
)

// IMessageReceiver defines the interface for consuming messages from a queue.
// Implementations of this interface handle the connection to message brokers
// and provide a channel for receiving incoming messages.
type IMessageReceiver[Obs any] interface {
	// ConsumeMessage returns a channel that yields incoming messages.
	// Each message includes the user's current state and the received message.
	// The channel should be closed when the adapter is stopped.
	ConsumeMessage() <-chan struct {
		UserState d_user.UserState[Obs]
		Message   d_message.Message
	}
}
