package output_router_api

import (
	"encoding/json"

	dto_user "github.com/irissonnlima/chatgraph-go/adapters/dto/user"
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

type Detail struct {
	Text string `json:"detail"`
}

type Message struct {
	TextDetail Detail `json:"text_detail"`
}

type TransferPayload struct {
	ChatID      dto_user.ChatID `json:"chat_id"`
	Menu        string          `json:"menu_id"`
	Message     Message         `json:"message"`
	UserMessage string          `json:"user_message,omitempty"`
}

func (r *RouterApi) TransferToMenu(chatID d_user.ChatID, transfer d_action.TransferToMenu, message d_message.Message) error {
	textMessage := message.EntireText()
	if transfer.UserMessage != "" {
		textMessage = transfer.UserMessage
	}

	payload := TransferPayload{
		ChatID: dto_user.ChatID{
			UserID:    chatID.UserID,
			CompanyID: chatID.CompanyID,
		},
		Menu: transfer.MenuID,
		Message: Message{
			TextDetail: Detail{
				Text: textMessage,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.post("/v1/actions/messages/transfer_to_menu", jsonPayload)
}
