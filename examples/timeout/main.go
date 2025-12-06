// Example: timeout - Demonstrates custom timeout configuration
package main

import (
	"log"
	"os"
	"time"

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

	// Timeout handler
	app.RegisterRoute("timeout_route", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		ctx.SendTextMessage("⏰ Operation timed out!")
		return ctx.NextRoute("start")
	})

	// Custom timeout handler for slow operations
	app.RegisterRoute("slow_timeout", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		ctx.SendTextMessage("⏰ The slow operation timed out!")
		return ctx.NextRoute("start")
	})

	app.RegisterRoute("loop_route", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		return &chatgraph.RedirectResponse{TargetRoute: "start"}
	})

	app.RegisterRoute("start", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		ctx.SendTextMessage("Type 'slow' to test a slow operation, or 'fast' for a fast one.")
		return ctx.NextRoute("handle_input")
	})

	app.RegisterRoute("handle_input", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		switch ctx.Message.EntireText() {
		case "slow":
			return &chatgraph.RedirectResponse{TargetRoute: "slow_operation"}
		case "fast":
			return &chatgraph.RedirectResponse{TargetRoute: "fast_operation"}
		default:
			ctx.SendTextMessage("Type 'slow' or 'fast'")
			return nil
		}
	})

	// Route with custom short timeout (5 seconds)
	app.RegisterRoute("slow_operation", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		ctx.SendTextMessage("Starting slow operation (will timeout in 5 seconds)...")

		// Simulate a slow operation
		time.Sleep(10 * time.Second)

		// This won't execute if timeout occurs
		ctx.SendTextMessage("Slow operation completed!")
		return ctx.NextRoute("start")
	}, chatgraph.RouterHandlerOptions{
		Timeout: &chatgraph.TimeoutRouteOps{
			Duration: 5 * time.Second,
			Route:    "slow_timeout",
		},
	})

	// Route with default timeout
	app.RegisterRoute("fast_operation", func(ctx *chatgraph.Context[Obs]) chatgraph.RouteReturn {
		ctx.SendTextMessage("Fast operation completed instantly! ⚡")
		return ctx.NextRoute("start")
	})

	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
