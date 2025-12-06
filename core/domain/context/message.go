package d_context

import d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"

// SendMessage sends a message to the specified chat ID.
// Returns an error if the message could not be sent.
func (c *ChatContext[Obs]) SendMessage(message d_message.Message) error {
	if c.Context.Err() != nil {
		return c.Context.Err()
	}

	return c.router.SendMessage(c.UserState.ChatID, message, c.UserState.Platform)
}

func (c *ChatContext[Obs]) SendTextMessage(text string) error {
	message := d_message.Message{
		TextMessage: d_message.TextMessage{
			Detail: text,
		},
	}

	return c.SendMessage(message)
}
