package service

import (
	"chatgraph/adapters/config"
	grpcclient "chatgraph/adapters/grpc"
	"chatgraph/adapters/queue/rabbitmq"
	domain_primitives "chatgraph/domain/primitives"
)

type RouteHandler struct {
	RouteFunc   func(*MessageContext) *RouteError
	OnErrorFunc func(*MessageContext, *RouteError)
}

type ChatbotApp struct {
	config     config.Config
	rabbitChan <-chan domain_primitives.UserCall

	routes map[string]RouteHandler

	grpc *grpcclient.Client
}

func NewChatbotApp() *ChatbotApp {

	config := config.LoadConfig()
	rabbit := rabbitmq.NewRabbitMQ(config)
	rabbitChan, err := rabbit.GetMessages()
	if err != nil {
		panic("Error getting messages: " + err.Error())
	}

	client, err := grpcclient.NewClient(config.GrpcURI)
	if err != nil {
		panic("Error connecting to grpc: " + err.Error())
	}

	return &ChatbotApp{
		config:     config,
		rabbitChan: rabbitChan,
		routes:     make(map[string]RouteHandler),
		grpc:       client,
	}
}

func (c *ChatbotApp) Start() {
	for ucall := range c.rabbitChan {
		route := ucall.UserState.Route.CurrentRoute()

		messageCtx := NewMessageContext(c.grpc, ucall.UserState, ucall.Message)

		routeHandler, ok := c.routes[route]
		if ok {
			go func(ctx *MessageContext) {
				err := routeHandler.RouteFunc(ctx)
				if err != nil {
					routeHandler.OnErrorFunc(ctx, err)
				}
			}(messageCtx)
		} else {
			panic("Route not found: " + route)
		}
	}
}

func (c *ChatbotApp) AddRoute(route string, handler RouteHandler) {
	if c.routes == nil {
		c.routes = make(map[string]RouteHandler)
	}
	c.routes[route] = handler
}
