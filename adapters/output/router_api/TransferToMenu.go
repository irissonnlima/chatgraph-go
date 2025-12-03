package output_router_api

import (
	d_action "chatgraph/core/domain/action"
	d_message "chatgraph/core/domain/message"
	d_user "chatgraph/core/domain/user"
)

func (r *RouterApi) TransferToMenu(chatID d_user.ChatID, transfer d_action.TransferToMenu, message d_message.Message) error {
	return nil
}
