// Package d_context provides the ChatContext type which encapsulates all
// information needed to process a chat message within a route handler.
package d_context

import (
	d_message "chatgraph/core/domain/message"
	d_user "chatgraph/core/domain/user"
	adapter_output "chatgraph/core/ports/adapters/output"
	"context"
	"time"
)

// ChatContext holds all the information needed to process a chat interaction.
// It embeds context.Context for cancellation and deadline propagation,
// and provides access to user state, the incoming message, and router services.
// It is generic over Obs, allowing custom observation data types.
type ChatContext[Obs any] struct {
	// Context is the standard Go context for cancellation and timeouts.
	context.Context
	// UserState contains the current state of the user's session.
	UserState d_user.UserState[Obs]
	// Message is the incoming message being processed.
	Message d_message.Message
	// router provides messaging and session management capabilities.
	router adapter_output.RouterService
}

// NewChatContext creates a new ChatContext with the provided parameters.
// This is the primary constructor for creating a ChatContext instance.
func NewChatContext[Obs any](
	userState d_user.UserState[Obs],
	message d_message.Message,

	router adapter_output.RouterService,

	timeout time.Duration,
) (ChatContext[Obs], context.CancelFunc) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	ctxChatbot := ChatContext[Obs]{
		Context:   ctx,
		UserState: userState,
		Message:   message,
		router:    router,
	}

	return ctxChatbot, cancel
}

// SendMessage sends a message to the specified chat ID.
// Returns an error if the message could not be sent.
func (c *ChatContext[Obs]) SendMessage(to d_user.ChatID, message d_message.Message) error {
	return c.router.SendMessage(to, message)
}
