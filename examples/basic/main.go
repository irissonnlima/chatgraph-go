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

	"github.com/irissonnlima/chatgraph-go/chatgraph"
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
	rabbit := chatgraph.NewRabbitMQ[Obs](
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_VHOST"),
		os.Getenv("RABBITMQ_QUEUE"),
	)

	// Create router API client
	routerApi := chatgraph.NewRouterApi(
		os.Getenv("ROUTER_API_URL"),
		os.Getenv("ROUTER_API_USER"),
		os.Getenv("ROUTER_API_PASSWORD"),
	)

	// Create the chatbot application
	app := chatgraph.NewApp(rabbit, routerApi)

	// Register timeout handler (required)
	app.RegisterRoute("timeout_route", handleTimeout)

	// Register loop handler (required)
	app.RegisterRoute("loop_route", handleLoop)

	// Register main routes
	app.RegisterRoute("start", handleStart)
	app.RegisterRoute("menu", handleMenu)

	// Start the application
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start the application: %v", err)
	}
}

func handleTimeout(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
	log.Println("Timeout route executed")
	ctx.SendTextMessage("Your request has timed out. Please try again.")
	return ctx.NextRoute("start")
}

func handleLoop(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
	log.Println("Loop detected:", ctx.UserState.Route.CurrentRepeated())
	ctx.SendTextMessage("Loop detected. Redirecting to start.")
	return &chatgraph.RedirectResponse{TargetRoute: "start"}
}

func handleStart(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
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

func handleMenu(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
	switch ctx.Message.EntireText() {
	case "end":
		return chatgraph.EndAction{ID: "session_ended"}
	case "start":
		return chatgraph.RedirectResponse{TargetRoute: "start"}
	default:
		ctx.SendTextMessage("Type 'end' to finish or 'start' to restart.")
		return nil
	}
}
