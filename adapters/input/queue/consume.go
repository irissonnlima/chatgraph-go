package rabbitmq

import (
	"encoding/json"
	"log"

	dto_message "github.com/irissonnlima/chatgraph-go/adapters/dto/message"
	dto_user "github.com/irissonnlima/chatgraph-go/adapters/dto/user"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

// MessageType represents the type of message received from the queue.
type MessageType string

const (
	// MessageTypeMessage indicates a regular chat message.
	MessageTypeMessage MessageType = "message"
	// MessageTypeStatus indicates a status update.
	MessageTypeStatus MessageType = "status"
	// MessageTypeEvent indicates an event notification.
	MessageTypeEvent MessageType = "event"
)

// QueueMessage represents the structure of a message received from RabbitMQ.
type QueueMessage struct {
	UserState dto_user.UserState  `json:"user_state"`
	Message   dto_message.Message `json:"message"`
}

// ConsumeMessage starts consuming messages from the RabbitMQ queue.
// It returns a channel that yields UserState and Message pairs for each
// message of type "message" received from the queue.
// The consumer runs in an infinite loop and automatically reconnects if the connection drops.
func (r *RabbitMQ[Obs]) ConsumeMessage() <-chan struct {
	UserState d_user.UserState[Obs]
	Message   d_message.Message
} {
	out := make(chan struct {
		UserState d_user.UserState[Obs]
		Message   d_message.Message
	})

	go func() {
		for {
			msgs, err := r.channel.Consume(
				r.queue, // queue name
				"",      // consumer tag
				true,    // auto-ack
				false,   // exclusive
				false,   // no-local
				false,   // no-wait
				nil,     // args
			)
			if err != nil {
				log.Printf("[RABBITMQ - ConsumeMessage] Error consuming queue: %v. Reconnecting...", err)
				r.reconnect(5)
				continue
			}

			log.Printf("[RABBITMQ - ConsumeMessage] Listening to queue: %s", r.queue)

			for msg := range msgs {
				var queueMsg QueueMessage

				if err := json.Unmarshal(msg.Body, &queueMsg); err != nil {
					log.Printf("[RABBITMQ - ConsumeMessage] Error unmarshalling message: %v", err)
					continue
				}

				userState := dto_user.UserStateToDomain[Obs](queueMsg.UserState)
				message := queueMsg.Message.ToDomain()

				out <- struct {
					UserState d_user.UserState[Obs]
					Message   d_message.Message
				}{
					UserState: userState,
					Message:   message,
				}
			}

			// If we get here, the channel was closed (connection lost)
			log.Printf("[RABBITMQ - ConsumeMessage] Channel closed. Reconnecting...")
			r.reconnect(5)
		}
	}()

	return out
}
