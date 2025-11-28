package rabbitmq

import (
	adapter_input "chatgraph/core/ports/adapters/input"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ[Obs any] struct {
	user     string
	password string
	host     string
	vhost    string
	queue    string

	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQ[Obs any](
	user string,
	password string,
	host string,
	vhost string,
	queue string,
) adapter_input.IMessageReceiver[Obs] {
	rabbit := RabbitMQ[Obs]{
		user:     user,
		password: password,
		host:     host,
		vhost:    vhost,
		queue:    queue,
	}

	rabbit.connect()

	return &rabbit
}

func (r *RabbitMQ[Obs]) connect() error {
	log.Println("[RABBITMQ - connect] - Attempting to connect to RabbitMQ...")

	if r.channel != nil {
		log.Println("[RABBITMQ - connect] - Closing existing channel before reconnecting...")
		_ = r.channel.Close()
	}

	if r.connection != nil {
		log.Println("[RABBITMQ - connect] - Closing existing connection before reconnecting...")
		_ = r.connection.Close()
	}

	var err error
	connStr := "amqp://" + r.user + ":" + r.password + "@" + r.host + "/" + r.vhost
	log.Println(fmt.Sprintf("[RABBITMQ - connect] - Connection string: %s", connStr))

	// Configure connection with heartbeat
	r.connection, err = amqp.DialConfig(connStr, amqp.Config{
		Heartbeat: 10 * time.Second, // heartbeat configuration
	})
	if err != nil {
		log.Println(fmt.Sprintf("[RABBITMQ - connect] - Error connecting to RabbitMQ: %s", connStr), err)
		return err
	}

	r.channel, err = r.connection.Channel()
	if err != nil {
		log.Println("[RABBITMQ - connect] - Error creating channel", err)
		return err
	}

	log.Println("[RABBITMQ - connect] - Connection established successfully to RabbitMQ.")
	return nil
}

func (r *RabbitMQ[Obs]) reconnect(timeTry uint) {
	log.Println("[RABBITMQ - reconnect] - Attempting to reconnect to RabbitMQ...")
	if timeTry == 0 {
		timeTry = 10
	}

	for i := range 5 {
		log.Println(fmt.Sprintf("[RABBITMQ - reconnect] %d/5: Attempting to reconnect to RabbitMQ...", i+1))
		err := r.connect()
		if err == nil {
			log.Println("[RABBITMQ - reconnect] Successfully reconnected to RabbitMQ!")
			break
		}

		log.Println(fmt.Sprintf("[RABBITMQ - reconnect] Error attempting to reconnect: %v. Trying again in %d seconds...", err, timeTry))
		time.Sleep(time.Duration(timeTry) * time.Second)
	}
}
