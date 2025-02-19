package main

import (
	d "chatgraph/domain/primitives"
	"chatgraph/service"

	dotenv "github.com/joho/godotenv"
)

func startRoute(ctx *service.MessageContext) (d.Router, *service.RouteError) {
	route := ctx.UserState.Route

	if ctx.Message.ContentMessage == "foo" {
		return route, &service.RouteError{
			TypeError: "Meu Tipo de Erro",
			Message:   "Mensagem de Erro",
		}
	}
	ctx.SendTextMessage("Olá, eu sou um chatbot. Como posso te ajudar?")
	return route.NextRoute(false, "start2"), nil
}

func routeError(ctx *service.MessageContext, err *service.RouteError) d.Router {
	ctx.SendTextMessage("Erro: " + err.Message)

	return ctx.UserState.Route.PreviousRoute(false)
}

func start2Route(ctx *service.MessageContext) (d.Router, *service.RouteError) {
	route := ctx.UserState.Route

	switch ctx.Message.ContentMessage {
	case "foo":
		ctx.SendTextMessage("Você digitou foo")
	case "bar":
		ctx.SendTextMessage("Você digitou bar")
	case "exit":
		ctx.SendTextMessage("encerrando conversa")
		ctx.EndChat("08fab84f-2ef9-4fcc-a504-c62e1475b938", "Observação de encerramento")
	case "error":
		return route.PreviousRoute(false), &service.RouteError{
			TypeError: "Meu Tipo de Erro",
			Message:   "Mensagem de Erro",
		}
	default:
		ctx.SendTextMessage("Você digitou algo diferente de foo e bar")
	}

	return route.NextRoute(true, "start"), nil
}

func main() {

	dotenv.Load()

	app := service.NewChatbotApp()

	app.AddRoute("start", service.RouteHandler{
		RouteFunc:   startRoute,
		OnErrorFunc: routeError,
	})

	app.AddRoute("start2", service.RouteHandler{
		RouteFunc:   start2Route,
		OnErrorFunc: routeError,
	})

	app.Start()
}
