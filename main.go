package main

import (
	"chatgraph/rabbitmq"
	"log"
	"os"

	"fmt"

	dotenv "github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PrintMessage(msg amqp.Delivery) {
	messageText := string(msg.Body)
	fmt.Println(messageText)
}

func main() {

	dotenv.Load()

	user := os.Getenv("RABBIT_USER")
	pass := os.Getenv("RABBIT_PASS")
	uri := os.Getenv("RABBIT_URI")
	vhost := os.Getenv("RABBIT_VHOST")
	queue := os.Getenv("RABBIT_QUEUE")

	log.Println("Iniciando o consumidor de mensagens...")
	log.Printf("RabbitUser:  %s\n", user)
	log.Printf("RabbitPass:  %s\n", pass)
	log.Printf("RabbitURI:   %s\n", uri)
	log.Printf("RabbitVhost: %s\n", vhost)
	log.Printf("RabbitQueue: %s\n", queue)

	rabbit := rabbitmq.NewRabbitMQ(
		user,
		pass,
		uri,
		vhost,
	)

	rabbit.ConsumeQueue(
		queue,
		PrintMessage,
	)
}
