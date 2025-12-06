package output_router_api

import (
	dto_action "github.com/irissonnlima/chatgraph-go/adapters/dto/action"
	dto_user "github.com/irissonnlima/chatgraph-go/adapters/dto/user"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
	"encoding/json"
)

type EndSessionPayload struct {
	ChatID   dto_user.ChatID      `json:"chat_id"`
	ActionID dto_action.EndAction `json:"end_action"`
}

func (r *RouterApi) EndSession(chatID d_user.ChatID, actionId string) error {
	payload := EndSessionPayload{
		ChatID: dto_user.ChatID{
			UserID:    chatID.UserID,
			CompanyID: chatID.CompanyID,
		},
		ActionID: dto_action.EndAction{
			ID: actionId,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.post("/v1/actions/session/end", jsonPayload)
}
