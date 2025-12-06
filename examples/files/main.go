// Example: files - Demonstrates file upload and sending
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/irissonnlima/chatgraph-go/chat"
)

type Obs struct{}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] No .env file found")
	}

	rabbit := chat.NewRabbitMQ[Obs](
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_VHOST"),
		os.Getenv("RABBITMQ_QUEUE"),
	)

	routerApi := chat.NewRouterApi(
		os.Getenv("ROUTER_API_URL"),
		os.Getenv("ROUTER_API_USER"),
		os.Getenv("ROUTER_API_PASSWORD"),
	)

	// Create the engine and register routes
	engine := chat.NewEngine[Obs]()

	engine.RegisterRoute("timeout_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendTextMessage("Timeout!")
		return ctx.NextRoute("start")
	})

	engine.RegisterRoute("loop_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		return &chat.RedirectResponse{TargetRoute: "start"}
	})

	engine.RegisterRoute("start", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendTextMessage("Send 'file' to receive a file, or 'upload' to upload from bytes.")
		return ctx.NextRoute("handle_input")
	})

	engine.RegisterRoute("handle_input", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		input := ctx.Message.EntireText()

		switch input {
		case "file":
			// Upload a file from disk and send it
			file, err := ctx.LoadFile("README.md")
			if err != nil {
				ctx.SendTextMessage("Error loading file: " + err.Error())
				return nil
			}
			if file != nil {
				ctx.SendMessage(chat.Message{
					TextMessage: chat.TextMessage{
						Detail: "Here's your file:",
					},
					File: *file,
				})
			}

		case "upload":
			// Create a file from bytes
			content := []byte("Hello, this is a test file content!")
			file, err := ctx.LoadFileBytes("test-file.txt", content)
			if err != nil {
				ctx.SendTextMessage("Error uploading file: " + err.Error())
				return nil
			}
			if file != nil {
				ctx.SendMessage(chat.Message{
					TextMessage: chat.TextMessage{
						Detail: "File created from bytes:",
					},
					File: *file,
				})
			}

		default:
			ctx.SendTextMessage("Unknown command. Try 'file' or 'upload'.")
		}

		return nil
	})

	// Create the app with the engine
	app := chat.NewApp(engine, rabbit, routerApi)

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
