package output_router_api

import (
	dto_file "chatgraph/adapters/dto/file"
	dto_message "chatgraph/adapters/dto/message"
	dto_user "chatgraph/adapters/dto/user"
	d_message "chatgraph/core/domain/message"
	d_user "chatgraph/core/domain/user"
	"encoding/json"
)

type SendMessagePayload struct {
	UserState dto_user.UserState  `json:"user_state"`
	Message   dto_message.Message `json:"message"`
}

func (r *RouterApi) SendMessage(to d_user.ChatID, message d_message.Message, platform string) error {

	if message.DisplayButton.IsEmpty() && len(message.Buttons) > 0 {
		message.DisplayButton = d_message.Button{
			Type:   d_message.POSTBACK,
			Title:  "Open",
			Detail: "Open buttons",
		}
	}

	buttons := make([]dto_message.Button, len(message.Buttons))
	for i, btn := range message.Buttons {
		buttons[i] = dto_message.Button{
			Type:   btn.Type.String(),
			Title:  btn.Title,
			Detail: btn.Detail,
		}
	}

	var displayButton *dto_message.Button = nil
	if !message.DisplayButton.IsEmpty() {
		displayButton = &dto_message.Button{
			Type:   message.DisplayButton.Type.String(),
			Title:  message.DisplayButton.Title,
			Detail: message.DisplayButton.Detail,
		}
	}

	var file *dto_file.File = nil
	if message.HasFile() {
		file = &dto_file.File{
			ID:   message.File.ID,
			Type: message.File.Type.String(),
			URL:  message.File.URL,
			Name: message.File.Name,
		}
	}

	payload := SendMessagePayload{
		UserState: dto_user.UserState{
			ChatID: &dto_user.ChatID{
				UserID:    to.UserID,
				CompanyID: to.CompanyID,
			},
			Platform: platform,
		},
		Message: dto_message.Message{
			TextMessage: &dto_message.TextMessage{
				Title:        message.TextMessage.Title,
				Detail:       message.TextMessage.Detail,
				Caption:      message.TextMessage.Caption,
				MentionedIds: message.TextMessage.MentionedIds,
			},
			Buttons:       buttons,
			DisplayButton: displayButton,
			File:          file,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return r.post("/v1/actions/messages/send", jsonPayload)
}
