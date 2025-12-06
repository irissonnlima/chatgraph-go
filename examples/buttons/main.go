// Example: buttons - Demonstrates sending messages with interactive buttons
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/irissonnlima/chatgraph-go/chatgraph"
)

type Obs struct{}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] No .env file found")
	}

	rabbit := chatgraph.NewRabbitMQ[Obs](
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_VHOST"),
		os.Getenv("RABBITMQ_QUEUE"),
	)

	routerApi := chatgraph.NewRouterApi(
		os.Getenv("ROUTER_API_URL"),
		os.Getenv("ROUTER_API_USER"),
		os.Getenv("ROUTER_API_PASSWORD"),
	)

	app := chatgraph.NewApp(rabbit, routerApi)

	app.RegisterRoute("timeout_route", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		ctx.SendTextMessage("Timeout! Please try again.")
		return ctx.NextRoute("start")
	})

	app.RegisterRoute("loop_route", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		return &chatgraph.RedirectResponse{TargetRoute: "start"}
	})

	app.RegisterRoute("start", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		// Send a message with buttons
		ctx.SendMessage(chatgraph.Message{
			TextMessage: chatgraph.TextMessage{
				Title:  "Welcome!",
				Detail: "Please choose an option:",
			},
			Buttons: []chatgraph.Button{
				{
					Type:   chatgraph.POSTBACK,
					Title:  "Option A",
					Detail: "option_a",
				},
				{
					Type:   chatgraph.POSTBACK,
					Title:  "Option B",
					Detail: "option_b",
				},
				{
					Type:   chatgraph.URL,
					Title:  "Visit Website",
					Detail: "https://example.com",
				},
			},
		})

		return ctx.NextRoute("handle_choice")
	})

	app.RegisterRoute("handle_choice", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		choice := ctx.Message.EntireText()

		switch choice {
		case "option_a":
			ctx.SendTextMessage("You selected Option A!")
		case "option_b":
			ctx.SendTextMessage("You selected Option B!")
		default:
			ctx.SendTextMessage("Unknown option: " + choice)
		}

		return &chatgraph.RedirectResponse{TargetRoute: "start"}
	})

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
