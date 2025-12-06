// Example: buttons - Demonstrates sending messages with interactive buttons
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

	app := chat.NewApp(rabbit, routerApi)

	app.RegisterRoute("timeout_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendTextMessage("Timeout! Please try again.")
		return ctx.NextRoute("start")
	})

	app.RegisterRoute("loop_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		return &chat.RedirectResponse{TargetRoute: "start"}
	})

	app.RegisterRoute("start", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		// Send a message with buttons
		ctx.SendMessage(chat.Message{
			TextMessage: chat.TextMessage{
				Title:  "Welcome!",
				Detail: "Please choose an option:",
			},
			Buttons: []chat.Button{
				{
					Type:   chat.POSTBACK,
					Title:  "Option A",
					Detail: "option_a",
				},
				{
					Type:   chat.POSTBACK,
					Title:  "Option B",
					Detail: "option_b",
				},
				{
					Type:   chat.URL,
					Title:  "Visit Website",
					Detail: "https://example.com",
				},
			},
		})

		return ctx.NextRoute("handle_choice")
	})

	app.RegisterRoute("handle_choice", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		choice := ctx.Message.EntireText()

		switch choice {
		case "option_a":
			ctx.SendTextMessage("You selected Option A!")
		case "option_b":
			ctx.SendTextMessage("You selected Option B!")
		default:
			ctx.SendTextMessage("Unknown option: " + choice)
		}

		return &chat.RedirectResponse{TargetRoute: "start"}
	})

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
