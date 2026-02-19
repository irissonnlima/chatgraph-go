# Chatgraph-Go

[![Go Reference](https://pkg.go.dev/badge/github.com/irissonnlima/chatgraph-go.svg)](https://pkg.go.dev/github.com/irissonnlima/chatgraph-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/irissonnlima/chatgraph-go)](https://goreportcard.com/report/github.com/irissonnlima/chatgraph-go)

A lightweight, flexible chatbot framework for Go with route-based conversation flow, timeout handling, and loop protection.

English | [Português](README.pt-br.md)

## Features

- 🚀 **Simple API** - Single import, intuitive route registration
- 🔄 **Route-based Flow** - Define conversation flows with named routes
- ⏱️ **Timeout Handling** - Automatic timeout with configurable duration and fallback routes
- 🔁 **Loop Protection** - Prevents infinite redirect loops automatically
- 📦 **Generic Observations** - Store custom session data with type safety
- 🔌 **Pluggable Adapters** - RabbitMQ input, REST API output (easily extensible)
- 📄 **File Support** - Upload, send, and download files with SHA256 deduplication
- ✅ **Well Tested** - Comprehensive test coverage for core packages

## Installation

```bash
go get github.com/irissonnlima/chatgraph-go/chat@latest
```

## Quick Start

```go
package main

import (
    "github.com/irissonnlima/chatgraph-go/chat"
)

// Define your observation type for session data
type Obs struct {
    Counter int `json:"counter"`
}

func main() {
    // Create adapters
    rabbit := chat.NewRabbitMQ[Obs]("user", "pass", "host", "vhost", "queue")
    router := chat.NewRouterApi("http://api-url", "user", "pass")
    engine := chat.NewEngine[Obs]()
    
    // Create app
    app := chat.NewApp(engine, rabbit, router)
    
    // Register routes
    engine.RegisterRoute("start", func(ctx *chat.Context[Obs]) chat.RouteReturn {
        ctx.SendTextMessage("Hello! Type something:")
        r := ctx.NextRoute("echo")
        return &r
    })
    
    engine.RegisterRoute("echo", func(ctx *chat.Context[Obs]) chat.RouteReturn {
        ctx.SendTextMessage("You said: " + ctx.Message.EntireText())
        r := ctx.NextRoute("start")
        return &r
    })
    
    // Required: timeout and loop handlers
    engine.RegisterRoute("timeout_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
        ctx.SendTextMessage("Session timed out!")
        r := ctx.NextRoute("start")
        return &r
    })
    
    engine.RegisterRoute("loop_route", func(ctx *chat.Context[Obs]) chat.RouteReturn {
        return &chat.RedirectResponse{TargetRoute: "start"}
    })
    
    app.Start()
}
```

## Architecture

Chatgraph follows a **hexagonal architecture** pattern, separating core domain logic from external adapters:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Application                              │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │  RabbitMQ   │──▶│  ChatbotApp │───▶│    Router API       │  │
│  │  (Input)    │    │  (Service)  │    │    (Output)         │  │
│  └─────────────┘    └──────┬──────┘    └─────────────────────┘  │
│                            │                                    │
│                     ┌──────▼──────┐                             │
│                     │   Routes    │                             │
│                     │  (Handlers) │                             │
│                     └─────────────┘                             │
└─────────────────────────────────────────────────────────────────┘
```

### Core Concepts

#### 1. Routes

Routes are named conversation states. Each route has a handler function that:

- Receives a `Context` with user state and message
- Sends messages to the user
- Returns the next action (next route, redirect, end session, etc.)

```go
app.RegisterRoute("greeting", func(ctx *chat.Context[Obs]) chat.RouteReturn {
    ctx.SendTextMessage("Hello!")
    return ctx.NextRoute("menu")  // User's next message goes to "menu" route
})
```

#### 2. Route Returns

Handlers can return different actions:

| Return Type | Behavior |
|-------------|----------|
| `ctx.NextRoute("name")` | Sets the route for the user's **next** message |
| `&RedirectResponse{TargetRoute: "name"}` | **Immediately** executes another route |
| `EndAction{ID: "reason"}` | Ends the conversation session |
| `TransferToMenu{MenuID: 1}` | Transfers user to a different menu |
| `nil` | Stays on the current route |

**NextRoute vs Redirect:**

```go
// NextRoute: Waits for user input, then executes "menu"
return ctx.NextRoute("menu")

// Redirect: Immediately executes "menu" without waiting
return &chat.RedirectResponse{TargetRoute: "menu"}
```

#### 3. Context

The `Context` provides access to:

```go
func handler(ctx *chat.Context[Obs]) chat.RouteReturn {
    // User information
    ctx.UserState.User.Name      // User's name
    ctx.UserState.ChatID         // Chat identifier
    ctx.UserState.Route          // Navigation history
    
    // Incoming message
    ctx.Message.EntireText()     // Full message text
    ctx.Message.TextMessage      // Structured text (Title, Detail, Footer)
    ctx.Message.Buttons          // Button responses
    ctx.Message.File             // File attachments
    
    // Session observations (custom data)
    obs := ctx.GetObservation()  // Get typed observation
    ctx.SetObservation(obs)      // Update observation
    
    // Send messages
    ctx.SendTextMessage("Hello!")
    ctx.SendMessage(chat.Message{...})
    
    // File operations
    ctx.LoadFile("path/to/file")           // Upload from disk
    ctx.LoadFileBytes("name.txt", []byte)  // Upload from bytes
    
    return ctx.NextRoute("next")
}
```

#### 4. Timeout Handling

Each route has a configurable timeout. When exceeded:

1. The handler execution is **cancelled** via context
2. User is **redirected** to the timeout route
3. No more messages are sent from the timed-out handler

```go
// Default: 5 minutes, redirects to "timeout_route"
app.RegisterRoute("slow_task", handler)

// Custom timeout: 30 seconds, redirects to "custom_timeout"
app.RegisterRoute("fast_task", handler, chat.RouterHandlerOptions{
    Timeout: &chat.TimeoutRouteOps{
        Duration: 30 * time.Second,
        Route:    "custom_timeout",
    },
})
```

**How it works internally:**

```
User Message ──▶ Handler Starts ──▶ [5 min timeout]
                      │
                      ├── Handler completes ──▶ Process result
                      │
                      └── Timeout exceeded ──▶ Cancel context
                                              └── Redirect to timeout_route
```

#### 5. Loop Protection

Prevents infinite redirect loops by counting consecutive visits to the same route:

```go
// Default: 3 consecutive visits, redirects to "loop_route"
// If route "A" redirects to "A" 3 times, user goes to "loop_route"
```

**How it works:**

```
A → A → A → A (4th time) ──▶ Redirect to loop_route
    │   │   │
    └───┴───┴── CurrentRepeated() = 3 > limit
```

#### 6. Messages with Buttons

Send interactive messages with clickable buttons:

```go
ctx.SendMessage(chat.Message{
    TextMessage: chat.TextMessage{
        Title:  "Choose an option",
        Detail: "Please select one:",
    },
    Buttons: []chat.Button{
        {Type: chat.POSTBACK, Title: "Option A", Detail: "option_a"},
        {Type: chat.POSTBACK, Title: "Option B", Detail: "option_b"},
        {Type: chat.URL, Title: "Visit Site", Detail: "https://example.com"},
    },
})
```

Button types:

- `POSTBACK`: Sends the `Detail` value back as user message
- `URL`: Opens the URL in user's browser

#### 7. Observations (Session Data)

Store custom typed data that persists across messages:

```go
type Obs struct {
    Step     int    `json:"step"`
    UserData string `json:"user_data"`
}

func handler(ctx *chat.Context[Obs]) chat.RouteReturn {
    obs := ctx.GetObservation()
    obs.Step++
    obs.UserData = ctx.Message.EntireText()
    ctx.SetObservation(obs)
    
    return ctx.NextRoute("next")
}
```

#### 8. File Handling

Upload and send files:

```go
// Upload from disk
file, err := ctx.LoadFile("document.pdf")
if err == nil && file != nil {
    ctx.SendMessage(chat.Message{File: *file})
}

// Upload from bytes (e.g., generated content)
content := []byte("Hello, World!")
file, err := ctx.LoadFileBytes("greeting.txt", content)
```

Files are deduplicated using SHA256 hash - uploading the same content twice returns the cached file.

**Download file content:**

```go
// Download file bytes from URL
if !ctx.Message.File.IsEmpty() {
    bytes, err := ctx.Message.File.Bytes()
    if err == nil {
        // Process the file bytes
        fmt.Printf("Downloaded %d bytes\n", len(bytes))
    }
}
```

## Configuration

### Default Options

```go
app := chat.NewApp(rabbit, router, chat.RouterHandlerOptions{
    Timeout: &chat.TimeoutRouteOps{
        Duration: 10 * time.Minute,  // Default timeout for all routes
        Route:    "timeout_route",
    },
    LoopCount: &chat.LoopCountRouteOps{
        Count: 5,                    // Allow 5 consecutive same-route visits
        Route: "loop_route",
    },
})
```

### Per-Route Options

```go
app.RegisterRoute("sensitive", handler, chat.RouterHandlerOptions{
    Timeout: &chat.TimeoutRouteOps{
        Duration: 1 * time.Minute,
        Route:    "sensitive_timeout",
    },
})
```

## Examples

See the [examples/](./examples/) directory for complete working examples:

- **basic/** - Simple chatbot with observations
- **buttons/** - Interactive buttons demo
- **files/** - File upload and download
- **timeout/** - Custom timeout configuration

## Project Structure

```
chatgraph-go/
├── chat/                # Unified public API package
│   └── chatgraph.go     # Type aliases and constructors
├── adapters/
│   ├── input/queue/     # RabbitMQ message consumer
│   └── output/router_api/  # REST API client
├── core/
│   ├── domain/          # Domain models
│   │   ├── action/      # Route return actions
│   │   ├── context/     # Chat context
│   │   ├── message/     # Message types
│   │   ├── route/       # Navigation history
│   │   ├── router/      # Handler options
│   │   └── user/        # User state
│   ├── ports/adapters/  # Adapter interfaces
│   └── service/         # Application service
└── examples/            # Usage examples
```

## Testing

Run tests with coverage:

```bash
go test ./... -cover

# Or use the coverage script for HTML report
./coverage.sh
# Open coverage/coverage.html in your browser
```

## License

MIT License - see [LICENSE](LICENSE) for details.
