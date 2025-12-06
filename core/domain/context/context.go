// Package d_context provides the ChatContext type which encapsulates all
// information needed to process a chat message within a route handler.
package d_context

import (
	"context"
	"time"

	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
	adapter_output "github.com/irissonnlima/chatgraph-go/core/ports/adapters/output"
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
	router adapter_output.IBotExecutor
}

// NewChatContext creates a new ChatContext with the provided parameters.
// This is the primary constructor for creating a ChatContext instance.
func NewChatContext[Obs any](
	userState d_user.UserState[Obs],
	message d_message.Message,

	router adapter_output.IBotExecutor,

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
