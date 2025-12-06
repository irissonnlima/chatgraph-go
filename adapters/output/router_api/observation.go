package output_router_api

import (
	dto_user "github.com/irissonnlima/chatgraph-go/adapters/dto/user"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
	"encoding/json"
)

type ObservationPayload struct {
	ChatId      dto_user.ChatID `json:"chat_id"`
	Observation string          `json:"observation"`
}

func (r *RouterApi) SetObservation(chatID d_user.ChatID, observation string) error {
	payload := ObservationPayload{
		ChatId: dto_user.ChatID{
			UserID:    chatID.UserID,
			CompanyID: chatID.CompanyID,
		},
		Observation: observation,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.post("/v1/actions/session/observation", jsonPayload)
}
