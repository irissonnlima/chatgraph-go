# Chatgraph Examples

This directory contains example applications demonstrating various features of the chatgraph-go framework.

## Setup

Before running any example, make sure you have a `.env` file in the example directory (or in the root of the project) with the following variables:

```env
# RabbitMQ Configuration
RABBITMQ_USER=your_user
RABBITMQ_PASSWORD=your_password
RABBITMQ_HOST=your_host
RABBITMQ_VHOST=your_vhost
RABBITMQ_QUEUE=your_queue

# Router API Configuration
ROUTER_API_URL=http://your-api-url
ROUTER_API_USER=your_api_user
ROUTER_API_PASSWORD=your_api_password
```

## Examples

### basic

A simple chatbot demonstrating:

- Basic route registration
- Using observations to store session data
- Sending text messages
- Route navigation

```bash
cd examples/basic
go run main.go
```

### buttons

Demonstrates interactive buttons:

- Sending messages with POSTBACK buttons
- Sending messages with URL buttons
- Handling button responses

```bash
cd examples/buttons
go run main.go
```

### files

Demonstrates file handling:

- Uploading files from disk
- Creating files from bytes
- Sending files in messages

```bash
cd examples/files
go run main.go
```

### timeout

Demonstrates timeout configuration:

- Default timeout handling
- Custom timeout per route
- Timeout redirection

```bash
cd examples/timeout
go run main.go
```

## Quick Start

The simplest way to create a chatbot:

```go
package main

import (
    "github.com/irissonnlima/chatgraph-go"
)

type Obs struct{} // Your observation type

func main() {
    rabbit := chat.NewRabbitMQ[Obs]("user", "pass", "host", "vhost", "queue")
    router := chat.NewRouterApi("url", "user", "pass")
    
    app := chat.NewApp(rabbit, router)
    
    app.RegisterRoute("start", func(ctx *chat.Context[Obs]) chat.RouteReturn {
        ctx.SendTextMessage("Hello World!")
        return nil
    })
    
    app.Start()
}
```
