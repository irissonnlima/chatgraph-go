package output_router_api

import (
	dto_user "chatgraph/adapters/dto/user"
	d_user "chatgraph/core/domain/user"
	"encoding/json"
)

type RoutePayload struct {
	ChatId dto_user.ChatID `json:"chat_id"`
	Route  string          `json:"route"`
}

func (r *RouterApi) SetRoute(chatID d_user.ChatID, route string) error {
	payload := RoutePayload{
		ChatId: dto_user.ChatID{
			UserID:    chatID.UserID,
			CompanyID: chatID.CompanyID,
		},
		Route: route,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.post("/v1/actions/session/route", jsonPayload)
}
