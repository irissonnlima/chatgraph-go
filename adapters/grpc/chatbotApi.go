package grpcclient

import (
	"chatgraph/adapters/grpc/chatbot"
	domain_primitives "chatgraph/domain/primitives"
	domain_response "chatgraph/domain/response"
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn

	UserStateService chatbot.UserStateServiceClient
	SendMessage      chatbot.SendMessageClient
	Transfer         chatbot.TransferClient
	EndChat          chatbot.EndChatClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:             conn,
		UserStateService: chatbot.NewUserStateServiceClient(conn),
		SendMessage:      chatbot.NewSendMessageClient(conn),
		Transfer:         chatbot.NewTransferClient(conn),
		EndChat:          chatbot.NewEndChatClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) InsertOrUpdateUserState(ctx context.Context, userState domain_primitives.UserState) error {

	state := &chatbot.UserState{
		ChatId: &chatbot.ChatID{
			UserId:    userState.ChatID.UserID,
			CompanyId: userState.ChatID.CompanyID,
		},
		Menu:        userState.Menu,
		Route:       userState.Route.HistoryRoute(),
		Observation: domain_primitives.NilString(userState.Observation),
	}
	status, err := c.UserStateService.InsertUpdateUserState(ctx, state)
	if err != nil {
		log.Println("Error requesting, inserting or updating user state: ", err)
		return err
	}

	if !status.Status {
		log.Println("Error status inserting or updating user state: ", status.Message)
		return errors.New(status.Message)
	}

	return nil
}

func (c *Client) DeleteUserState(ctx context.Context, chatID *chatbot.ChatID) (*chatbot.RequestStatus, error) {
	return c.UserStateService.DeleteUserState(ctx, chatID)
}

func (c *Client) GetUserState(ctx context.Context, chatID *chatbot.ChatID) (*chatbot.UserState, error) {
	return c.UserStateService.GetUserState(ctx, chatID)
}

func (c *Client) GetAllUserStates(ctx context.Context) (*chatbot.UserStateList, error) {
	return c.UserStateService.GetAllUserStates(ctx, &chatbot.Void{})
}

func (c *Client) SendMessageMsg(
	ctx context.Context,
	chatID domain_primitives.ChatID,
	message domain_response.ResponseMessage,
) error {

	buttons := make([]*chatbot.Button, 0, len(message.Buttons))
	for _, b := range message.Buttons {
		var typeBtn string
		switch b.Type {
		case domain_response.Postback:
			typeBtn = "postback"
		case domain_response.URL:
			typeBtn = "url"
		default:
			typeBtn = "postback"
		}

		buttons = append(buttons, &chatbot.Button{
			Type:   typeBtn,
			Title:  b.Title,
			Detail: b.Detail,
		})
	}

	var displayButton *chatbot.Button
	if len(buttons) > 0 {
		if message.DiplayButton.Title == "" {
			message.DiplayButton.Title = "Ver mais"
		}

		displayButton = &chatbot.Button{
			Type:   "postback",
			Title:  message.DiplayButton.Title,
			Detail: message.DiplayButton.Detail,
		}
	}

	msg := &chatbot.Message{
		ChatId: &chatbot.ChatID{
			UserId:    chatID.UserID,
			CompanyId: chatID.CompanyID,
		},
		Message: &chatbot.TextMessage{
			Type:    message.TextMessage.Type,
			Title:   message.TextMessage.Title,
			Detail:  message.TextMessage.Detail,
			Caption: message.TextMessage.Caption,
		},
		Buttons:       buttons,
		DisplayButton: displayButton,
	}

	result, err := c.SendMessage.SendMessage(ctx, msg)
	if err != nil {
		log.Println("Error sending message to user: ", err)
		return err
	}

	if !result.Status {
		log.Println("Error status sending message to user: ", result.Message)
		return errors.New(result.Message)
	}

	return nil
}

func (c *Client) GetAllCampaigns(ctx context.Context) (*chatbot.CampaignsList, error) {
	return c.Transfer.GetAllCampaigns(ctx, &chatbot.Void{})
}

func (c *Client) GetCampaignID(ctx context.Context, name *chatbot.CampaignName) (*chatbot.CampaignDetails, error) {
	return c.Transfer.GetCampaignID(ctx, name)
}

// TransferToHuman chama o RPC TransferToHuman do serviço Transfer.
func (c *Client) TransferToHuman(ctx context.Context, req *chatbot.TransferToHumanRequest) (*chatbot.RequestStatus, error) {
	return c.Transfer.TransferToHuman(ctx, req)
}

// TransferToMenu chama o RPC TransferToMenu do serviço Transfer.
func (c *Client) TransferToMenu(ctx context.Context, req *chatbot.TransferToMenuRequest) (*chatbot.RequestStatus, error) {
	return c.Transfer.TransferToMenu(ctx, req)
}

// GetAllTabulations chama o RPC GetAllTabulations do serviço EndChat.
func (c *Client) GetAllTabulations(ctx context.Context) (*chatbot.TabulationsList, error) {
	return c.EndChat.GetAllTabulations(ctx, &chatbot.Void{})
}

// GetTabulationID chama o RPC GetTabulationID do serviço EndChat.
func (c *Client) GetTabulationID(ctx context.Context, name *chatbot.TabulationName) (*chatbot.TabulationDetails, error) {
	return c.EndChat.GetTabulationID(ctx, name)
}

// EndChat chama o RPC EndChat do serviço EndChat.
func (c *Client) EndChatService(ctx context.Context, req *chatbot.EndChatRequest) (*chatbot.RequestStatus, error) {
	return c.EndChat.EndChat(ctx, req)
}
