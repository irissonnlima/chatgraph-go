package rabbitmq

import (
	"chatgraph/adapters/queue"
	domain_primitives "chatgraph/domain/primitives"
	"encoding/json"
	"log"
)

type QueueOptions struct {
	Durable            bool
	AutoDelete         bool
	Exclusive          bool
	NoWait             bool
	TTL                int
	ExpiresSeconds     int
	DeadLetterExchange string
}

func (rabbit *RabbitMQ) GetMessages() (<-chan domain_primitives.UserCall, error) {
	if rabbit.channel == nil || rabbit.connection == nil || rabbit.channel.IsClosed() {
		log.Println("🐰 Reconnectando ao RabbitMQ")
		rabbit.reconnect(10)
	}

	log.Println("🐰 Consumindo mensagens da fila:", rabbit.Queue)
	msgs, err := rabbit.channel.Consume(
		rabbit.Queue, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Println("❌ Erro ao consumir mensagem:", err)
		return nil, err
	}

	out := make(chan domain_primitives.UserCall)

	go func() {
		defer close(out)
		for d := range msgs {
			var messageJson queue.MessageJson
			if err := json.Unmarshal(d.Body, &messageJson); err != nil {
				log.Println("❌ Erro ao parsear mensagem JSON:", err)
				continue
			}

			messageDomain := domain_primitives.UserCall{
				UserState: domain_primitives.UserState{
					ChatID: domain_primitives.ChatID{
						UserID:    messageJson.UserState.ChatID.UserID,
						CompanyID: messageJson.UserState.ChatID.CompanyID,
					},
					Menu:        messageJson.UserState.Menu,
					Route:       messageJson.UserState.Route,
					Observation: messageJson.UserState.Observation,
					Protocol:    messageJson.UserState.Protocol,
				},
				Message: domain_primitives.Message{
					TypeMessage:    messageJson.TypeMessage,
					ContentMessage: messageJson.ContentMessage,
				},
			}

			out <- messageDomain
		}
	}()

	return out, nil
}
