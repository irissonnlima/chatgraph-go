package main

import (
	"chatgraph/service"

	dotenv "github.com/joho/godotenv"
)

func startRoute(ctx *service.MessageContext) *service.RouteError {
	if ctx.Message.ContentMessage == "foo" {
		return &service.RouteError{
			TypeError: "Meu Tipo de Erro",
			Message:   "Mensagem de Erro",
		}
	}
	ctx.SendTextMessage("start", "Olá, eu sou um chatbot. Como posso te ajudar?")
	return nil
}

func startRouteError(ctx *service.MessageContext, err *service.RouteError) {
	ctx.SendTextMessage("start", "Erro: "+err.Message)
}

func main() {

	dotenv.Load()

	app := service.NewChatbotApp()

	app.AddRoute("start", service.RouteHandler{
		RouteFunc:   startRoute,
		OnErrorFunc: startRouteError,
	})

	app.Start()
}
