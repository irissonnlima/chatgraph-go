package output_router_api

import (
	d_action "github.com/irissonnlima/chatgraph-go/core/domain/action"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
)

func (r *RouterApi) TransferToMenu(chatID d_user.ChatID, transfer d_action.TransferToMenu, message d_message.Message) error {
	return nil
}
