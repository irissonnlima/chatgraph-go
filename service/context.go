package service

import (
	domain_primitives "chatgraph/domain/primitives"
	domain_response "chatgraph/domain/response"
	"context"
	"log"
)

type MessageContext struct {
	context.Context
	Route     domain_primitives.Router
	UserState domain_primitives.UserState
	Message   domain_primitives.Message
}

func (ctx *MessageContext) SendMessage(nextRoute string, message domain_response.ResponseMessage) {
	log.Println("Sending message: ", message)
}

func (ctx *MessageContext) SendTextMessage(nextRoute string, message string) {
	log.Println("Sending TextMessage: ", message)
}

func (ctx *MessageContext) EndChat(tabulationID string, observation string) {
	log.Println("Ending chat", tabulationID, " obs: ", observation)
}

func (ctx *MessageContext) TransferToHuman(campaignID string, observation string) {
	log.Println("Transfering chat", campaignID, " obs: ", observation)
}

func (ctx *MessageContext) TransferToMenu(menu string, observation string) {
	log.Println("Transfering to id_menu: ", menu, " obs: ", observation)
}
