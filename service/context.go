package service

import (
	grpcclient "chatgraph/adapters/grpc"
	domain_primitives "chatgraph/domain/primitives"
	domain_response "chatgraph/domain/response"
	"context"
	"log"
)

type MessageContext struct {
	context.Context

	grpcclient *grpcclient.Client

	Route     domain_primitives.Router
	UserState domain_primitives.UserState
	Message   domain_primitives.Message
}

func NewMessageContext(
	grpcclient *grpcclient.Client,
	userState domain_primitives.UserState,
	message domain_primitives.Message,
) *MessageContext {

	return &MessageContext{
		Context:    context.Background(),
		grpcclient: grpcclient,
		Route:      userState.Route,
		UserState:  userState,
		Message:    message,
	}
}

func (ctx *MessageContext) updateRoute(nextRoute domain_primitives.Router) {

	updateUstate := domain_primitives.UserState{
		ChatID: ctx.UserState.ChatID,
		Menu:   ctx.UserState.Menu,
		Route:  nextRoute,
	}
	err := ctx.grpcclient.InsertOrUpdateUserState(ctx, updateUstate)
	if err != nil {
		log.Println("Error!! updating route: ", err)
		return
	}

	ctx.UserState.Route = nextRoute
	ctx.Route = nextRoute
}

func (ctx *MessageContext) SendMessage(message domain_response.ResponseMessage) {
	log.Println("Sending message: ", message)

	err := ctx.grpcclient.SendMessageMsg(ctx, ctx.UserState.ChatID, message)
	if err != nil {
		log.Println("Error!! sending message: ", err)
	}
}

func (ctx *MessageContext) SendTextMessage(message string) {
	log.Println("Sending TextMessage: ", message)

	err := ctx.grpcclient.SendMessageMsg(
		ctx,
		ctx.UserState.ChatID,
		domain_response.ResponseMessage{
			TextMessage: domain_response.TextMessage{
				Type:    "message",
				Title:   "",
				Detail:  message,
				Caption: "",
			},
			Buttons:      nil,
			DiplayButton: domain_response.Button{},
		},
	)
	if err != nil {
		log.Println("Error!! sending message: ", err)
	}
}

func (ctx *MessageContext) EndChat(tabulationID string, observation string) {
	log.Println("Ending chat", tabulationID, " obs: ", observation)

	err := ctx.grpcclient.EndChatService(ctx, ctx.UserState.ChatID, tabulationID, observation)
	if err != nil {
		log.Println("Error!! ending chat: ", err)
	}
}

func (ctx *MessageContext) TransferToHuman(campaignID string, observation string) {
	log.Println("Transfering chat", campaignID, " obs: ", observation)
}

func (ctx *MessageContext) TransferToMenu(menu string, observation string) {
	log.Println("Transfering to id_menu: ", menu, " obs: ", observation)
}
