// Example: timeout - Demonstrates custom timeout configuration
package main

import (
	"log"
	"os"
	"time"

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

	// Timeout handler
	engine.RegisterRoute("timeout_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendTextMessage("⏰ Operation timed out!")
		return ctx.NextRoute("start")
	})

	// Custom timeout handler for slow operations
	engine.RegisterRoute("slow_timeout", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendTextMessage("⏰ The slow operation timed out!")
		return ctx.NextRoute("start")
	})

	engine.RegisterRoute("loop_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		return &chat.RedirectResponse{TargetRoute: "start"}
	})

	engine.RegisterRoute("start", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendTextMessage("Type 'slow' to test a slow operation, or 'fast' for a fast one.")
		return ctx.NextRoute("handle_input")
	})

	engine.RegisterRoute("handle_input", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		switch ctx.Message.EntireText() {
		case "slow":
			return &chat.RedirectResponse{TargetRoute: "slow_operation"}
		case "fast":
			return &chat.RedirectResponse{TargetRoute: "fast_operation"}
		default:
			ctx.SendTextMessage("Type 'slow' or 'fast'")
			return nil
		}
	})

	// Route with custom short timeout (5 seconds)
	engine.RegisterRoute("slow_operation", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendTextMessage("Starting slow operation (will timeout in 5 seconds)...")

		// Simulate a slow operation
		time.Sleep(10 * time.Second)

		// This won't execute if timeout occurs
		ctx.SendTextMessage("Slow operation completed!")
		return ctx.NextRoute("start")
	}, chat.RouterHandlerOptions{
		Timeout: &chat.TimeoutRouteOps{
			Duration: 5 * time.Second,
			Route:    "slow_timeout",
		},
	})

	// Route with default timeout
	engine.RegisterRoute("fast_operation", func(ctx *chat.Context[Obs]) chat.RouteReturn {
		ctx.SendTextMessage("Fast operation completed instantly! ⚡")
		return ctx.NextRoute("start")
	})

	// Create the app with the engine
	app := chat.NewApp(engine, rabbit, routerApi)

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
