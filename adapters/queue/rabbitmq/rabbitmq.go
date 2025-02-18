package rabbitmq

import (
	"chatgraph/adapters/config"
	"chatgraph/adapters/queue"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	user       string
	password   string
	host       string
	vhost      string
	Queue      string
	connection *amqp.Connection
	channel    *amqp.Channel

	queue.IQueue
}

func NewRabbitMQ(cfg config.Config) *RabbitMQ {
	rabbit := RabbitMQ{
		Queue:    cfg.RabbitMQQueue,
		user:     cfg.RabbitMQUser,
		password: cfg.RabbitMQPassword,
		host:     cfg.RabbitMQHost,
		vhost:    cfg.RabbitMQVHost,
	}

	rabbit.connect()

	return &rabbit
}

func (rabbit *RabbitMQ) connect() error {
	if rabbit.channel != nil {
		_ = rabbit.channel.Close()
	}
	if rabbit.connection != nil {
		_ = rabbit.connection.Close()
	}

	var err error
	connStr := "amqp://" + rabbit.user + ":" + rabbit.password + "@" + rabbit.host + "/" + rabbit.vhost

	// Configurar a conexão com heartbeat
	rabbit.connection, err = amqp.DialConfig(connStr, amqp.Config{
		Heartbeat: 10 * time.Second, // Configuração do heartbeat
	})
	if err != nil {
		return err
	}

	rabbit.channel, err = rabbit.connection.Channel()
	if err != nil {
		return err
	}

	log.Println("Conexão estabelecida com sucesso ao RabbitMQ.")
	return nil
}

func (rabbit *RabbitMQ) reconnect(timeTry uint) {
	if timeTry == 0 {
		timeTry = 10
	}

	for i := range 5 {
		log.Printf("%d/5: Tentando reconectar ao RabbitMQ...\n", i+1)
		err := rabbit.connect()
		if err == nil {
			log.Println("Reconectado ao RabbitMQ com sucesso!")
			break
		}

		log.Printf("Erro ao tentar reconectar: %v. Tentando novamente em %d segundos...", err, timeTry)
		time.Sleep(time.Duration(timeTry) * time.Second)
	}
}
