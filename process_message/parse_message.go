package process_message

import (
	"chatgraph/models"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ParseMessage(msg amqp.Delivery) (error, *models.JsonMessage) {
	var message models.JsonMessage
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Printf("Erro ao fazer parse da mensagem recebida: %s\n", err)
		return err, nil
	}
	return nil, &message
}
