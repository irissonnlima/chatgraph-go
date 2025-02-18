package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GrpcURI string

	RabbitMQUser               string
	RabbitMQPassword           string
	RabbitMQHost               string
	RabbitMQVHost              string
	RabbitMQTTL                int
	RabbitMQExpires            int
	RabbitMQDeadLetterExchange string
	RabbitMQQueue              string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to environment variables.")
	}

	return Config{
		GrpcURI: os.Getenv("GRPC_URI"),

		RabbitMQQueue:    os.Getenv("RABBITMQ_QUEUE"),
		RabbitMQUser:     os.Getenv("RABBITMQ_USER"),
		RabbitMQPassword: os.Getenv("RABBITMQ_PASSWORD"),
		RabbitMQHost:     os.Getenv("RABBITMQ_HOST"),
		RabbitMQVHost:    os.Getenv("RABBITMQ_VHOST"),
	}
}
