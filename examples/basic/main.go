// Example: basic - A simple chatbot with basic routing
//
// This example demonstrates:
// - Setting up a chatbot with RabbitMQ and Router API
// - Registering routes with handlers
// - Using observations to store session data
// - Sending messages and handling user input
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/irissonnlima/chatgraph-go/chat"
)

// Obs defines the custom observation data for this chatbot.
// This data persists across messages within a session.
type Obs struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] No .env file found, using environment variables")
	}

	// Create message receiver (RabbitMQ)
	rabbit := chat.NewRabbitMQ[Obs](
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_VHOST"),
		os.Getenv("RABBITMQ_QUEUE"),
	)

	// Create router API client
	routerApi := chat.NewRouterApi(
		os.Getenv("ROUTER_API_URL"),
		os.Getenv("ROUTER_API_USER"),
		os.Getenv("ROUTER_API_PASSWORD"),
	)

	// Create the engine and register routes
	engine := chat.NewEngine[Obs]()
	engine.RegisterRoute("timeout_route", handleTimeout)
	engine.RegisterRoute("loop_route", handleLoop)
	engine.RegisterRoute("start", handleStart)
	engine.RegisterRoute("menu", handleMenu)

	// Create the chatbot application with the engine
	app := chat.NewApp(engine, rabbit, routerApi)

	// Start the application
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start the application: %v", err)
	}
}

func handleTimeout(ctx *chat.Context[Obs]) chat.RouteReturn {
	log.Println("Timeout route executed")
	ctx.SendTextMessage("Your request has timed out. Please try again.")
	return ctx.NextRoute("start")
}

func handleLoop(ctx *chat.Context[Obs]) chat.RouteReturn {
	log.Println("Loop detected:", ctx.UserState.Route.CurrentRepeated())
	ctx.SendTextMessage("Loop detected. Redirecting to start.")
	return &chat.RedirectResponse{TargetRoute: "start"}
}

func handleStart(ctx *chat.Context[Obs]) chat.RouteReturn {
	log.Println("Start route executed")
	log.Printf("User: %s", ctx.UserState.User.Name)

	// Get and update observation
	obs := ctx.GetObservation()
	ctx.SendTextMessage(fmt.Sprintf("Hello %s!", ctx.UserState.User.Name))
	ctx.SendTextMessage(fmt.Sprintf("Field1: %s, Field2: %d", obs.Field1, obs.Field2))

	// Toggle field1 value
	if obs.Field1 == "updated" {
		obs.Field1 = "changed"
	} else {
		obs.Field1 = "updated"
	}
	obs.Field2++
	ctx.SetObservation(obs)

	return ctx.NextRoute("menu")
}

func handleMenu(ctx *chat.Context[Obs]) chat.RouteReturn {
	switch ctx.Message.EntireText() {
	case "end":
		return chat.EndAction{ID: "session_ended"}
	case "start":
		return chat.RedirectResponse{TargetRoute: "start"}
	default:
		ctx.SendTextMessage("Type 'end' to finish or 'start' to restart.")
		return nil
	}
}
