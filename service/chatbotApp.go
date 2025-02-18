package service

import (
	"chatgraph/adapters/config"
	"chatgraph/adapters/queue/rabbitmq"
	domain_primitives "chatgraph/domain/primitives"
	"context"
)

type RouteHandler struct {
	RouteFunc   func(*MessageContext) *RouteError
	OnErrorFunc func(*MessageContext, *RouteError)
}

type ChatbotApp struct {
	config     config.Config
	rabbitChan <-chan domain_primitives.UserCall

	routes map[string]RouteHandler
}

func NewChatbotApp() *ChatbotApp {

	config := config.LoadConfig()
	rabbit := rabbitmq.NewRabbitMQ(config)
	rabbitChan, err := rabbit.GetMessages()
	if err != nil {
		panic("Error getting messages: " + err.Error())
	}

	return &ChatbotApp{
		config:     config,
		rabbitChan: rabbitChan,
	}
}

func (c *ChatbotApp) Start() {
	for ucall := range c.rabbitChan {
		route := ucall.UserState.Route
		if route == nil {
			route = new(string)
			*route = "start"
			ucall.UserState.Route = route
		}

		messageCtx := &MessageContext{
			Context:   context.Background(),
			Route:     domain_primitives.NewRouter(*route),
			UserState: ucall.UserState,
			Message:   ucall.Message,
		}

		routeHandler, ok := c.routes[*route]
		if ok {
			go func(ctx *MessageContext) {
				err := routeHandler.RouteFunc(ctx)
				if err != nil {
					routeHandler.OnErrorFunc(ctx, err)
				}
			}(messageCtx)
		} else {
			panic("Route not found: " + *route)
		}
	}
}

func (c *ChatbotApp) AddRoute(route string, handler RouteHandler) {
	if c.routes == nil {
		c.routes = make(map[string]RouteHandler)
	}
	c.routes[route] = handler
}
