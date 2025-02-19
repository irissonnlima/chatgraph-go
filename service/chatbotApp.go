package service

import (
	"chatgraph/adapters/config"
	grpcclient "chatgraph/adapters/grpc"
	"chatgraph/adapters/queue/rabbitmq"
	domain_primitives "chatgraph/domain/primitives"
)

type RouteHandler struct {
	RouteFunc   func(*MessageContext) (domain_primitives.Router, *RouteError)
	OnErrorFunc func(*MessageContext, *RouteError) domain_primitives.Router
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

func ProcessRoute(ctx *MessageContext, routerMap map[string]RouteHandler, router domain_primitives.Router) {
	route := router.CurrentRoute()
	if route == "error" {
		prevRouter := router.PreviousRoute(router.IsRedirect())
		route = prevRouter.CurrentRoute()
	}

	routeHandler, ok := routerMap[route]
	if ok {
		route, err := routeHandler.RouteFunc(ctx)
		if err != nil {
			ctx.updateRoute(ctx.Route.NextRoute(false, "error"))
			route = routeHandler.OnErrorFunc(ctx, err)
		}

		ctx.updateRoute(route)
		if route.IsRedirect() {
			ProcessRoute(ctx, routerMap, route)
		}
	} else {
		panic("Route not found: " + route)
	}
}

func (c *ChatbotApp) Start() {
	for ucall := range c.rabbitChan {
		messageCtx := NewMessageContext(c.grpc, ucall.UserState, ucall.Message)
		go ProcessRoute(messageCtx, c.routes, ucall.UserState.Route)
	}
}

func (c *ChatbotApp) AddRoute(route string, handler RouteHandler) {
	if c.routes == nil {
		c.routes = make(map[string]RouteHandler)
	}
	c.routes[route] = handler
}
