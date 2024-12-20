package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueOptions struct {
	Durable            bool
	AutoDelete         bool
	Exclusive          bool
	NoWait             bool
	TTL                int
	DeadLetterExchange string
}

func (rabbit RabbitMQ) CreateQueue(queue string, routingKey string, exchange string, options QueueOptions) (q amqp.Queue, err error) {

	args := amqp.Table{}
	if options.TTL > 0 {
		args["x-message-ttl"] = options.TTL
	}
	if options.DeadLetterExchange != "" {
		args["x-dead-letter-exchange"] = options.DeadLetterExchange
	}

	q, errLocal := rabbit.channel.QueueDeclare(
		queue,
		options.Durable,
		options.AutoDelete,
		options.Exclusive,
		options.NoWait,
		args,
	)
	if errLocal != nil {
		return q, errLocal
	}
	errLocal = rabbit.channel.QueueBind(
		queue,      // nome da fila
		routingKey, // routing key (chave de roteamento)
		exchange,   // nome da exchange
		false,      // no-wait
		nil,        // argumentos adicionais
	)
	if errLocal != nil {
		return q, errLocal
	}
	return q, nil
}

func (rabbit *RabbitMQ) ConsumeQueue(queue string, processMessage func(amqp.Delivery)) {
	msgs, err := rabbit.channel.Consume(
		queue,
		"",    // Nome do consumidor (deixe vazio para que o RabbitMQ gere um nome)
		true,  // Auto-Ack (true para confirmar automaticamente as mensagens)
		false, // Exclusivo
		false, // No-local (não aplicável)
		false, // No-wait
		nil,   // Argumentos adicionais
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			processMessage(d)
		}

		log.Println("Conexão perdida com RabbitMQ. Tentando reconectar...")
		rabbit.reconnect(10)
	}()

	log.Printf("Esperando mensagens...")
	<-forever // Mantém o processo rodando
}

func (rabbit *RabbitMQ) PublishMessage(routingKey string, exchange string, message string) error {
	if rabbit.channel == nil || rabbit.connection == nil || rabbit.channel.IsClosed() {
		rabbit.reconnect(10)
	}

	err := rabbit.channel.Publish(
		exchange,   // Exchange
		routingKey, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		return err
	}

	return nil
}

func (rabbit *RabbitMQ) PublishMessageWithReply(routingKey string, exchange string, message string) (response string, err error) {

	if rabbit.channel == nil || rabbit.connection == nil || rabbit.channel.IsClosed() {
		rabbit.reconnect(10)
	}

	timeout := 60
	correlationID := uuid.New().String()

	replyRoutingKey := routingKey + ".reply." + correlationID
	fmt.Printf("Reply Routing Key: %s\n", replyRoutingKey)

	_, erroDeclareQueue := rabbit.channel.QueueDeclare(
		replyRoutingKey,
		false, // durable
		true,  // autoDelete
		true,
		false,
		nil,
	)
	if erroDeclareQueue != nil {
		return "", erroDeclareQueue
	}

	msgs, err := rabbit.channel.Consume(
		replyRoutingKey, // Fila para ouvir a resposta
		"",              // Nome do consumidor
		true,            // Auto-Ack (true para confirmar automaticamente as mensagens)
		false,           // Exclusivo
		false,           // No-local
		false,           // No-wait
		nil,             // Argumentos adicionais
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer for reply: %v", err)
		return "", err
	}

	// Publicar a mensagem com o campo ReplyTo e o correlation ID
	err = rabbit.channel.Publish(
		exchange,   // Exchange
		routingKey, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(message),
			ReplyTo:       replyRoutingKey, // Nome da fila para resposta
			CorrelationId: correlationID,   // Correlation ID único para rastrear a resposta
		})
	if err != nil {
		log.Fatalf("Failed to publish a message with reply: %v", err)
		return "", err
	}

	// Canal para receber a resposta
	replyChan := make(chan string)

	go func() {
		for d := range msgs {
			// Verificar se a mensagem de resposta tem o mesmo Correlation ID
			if d.CorrelationId == correlationID {
				replyChan <- string(d.Body)
				break
			}
		}
	}()

	// Espera pela resposta com timeout
	select {
	case reply := <-replyChan:
		_, err := rabbit.channel.QueueDelete(replyRoutingKey, false, false, false)
		if err != nil {
			log.Printf("Failed to delete reply queue: %v", err)
		} else {
			log.Printf("Reply queue %s deleted after %d seconds\n", replyRoutingKey, timeout)
		}
		return reply, nil

	case <-time.After(time.Duration(timeout) * time.Second): // Timeout de 10 segundos
		log.Println("Timeout")
		_, err := rabbit.channel.QueueDelete(replyRoutingKey, false, false, false)
		if err != nil {
			log.Printf("TIMEOUT: Failed to delete reply queue: %v", err)
		} else {
			log.Printf("TIMEOUT: Reply queue %s deleted after %d seconds\n", replyRoutingKey, timeout)
		}
		return "", errors.New("Timeout")
	}
}
