package dto_message

import (
	"testing"

	dto_file "github.com/irissonnlima/chatgraph-go/adapters/dto/file"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
)

func TestTextMessage_ToDomain(t *testing.T) {
	tests := []struct {
		name string
		tm   TextMessage
	}{
		{
			name: "full text message",
			tm: TextMessage{
				ID:           "123",
				Title:        "Title",
				Detail:       "Detail",
				Caption:      "Caption",
				MentionedIds: []string{"user1", "user2"},
			},
		},
		{
			name: "empty text message",
			tm:   TextMessage{},
		},
		{
			name: "partial text message",
			tm: TextMessage{
				Title:  "Only Title",
				Detail: "Only Detail",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tm.ToDomain()
			if got.ID != tt.tm.ID {
				t.Errorf("TextMessage.ToDomain().ID = %v, want %v", got.ID, tt.tm.ID)
			}
			if got.Title != tt.tm.Title {
				t.Errorf("TextMessage.ToDomain().Title = %v, want %v", got.Title, tt.tm.Title)
			}
			if got.Detail != tt.tm.Detail {
				t.Errorf("TextMessage.ToDomain().Detail = %v, want %v", got.Detail, tt.tm.Detail)
			}
			if got.Caption != tt.tm.Caption {
				t.Errorf("TextMessage.ToDomain().Caption = %v, want %v", got.Caption, tt.tm.Caption)
			}
			if len(got.MentionedIds) != len(tt.tm.MentionedIds) {
				t.Errorf("TextMessage.ToDomain().MentionedIds length = %v, want %v", len(got.MentionedIds), len(tt.tm.MentionedIds))
			}
		})
	}
}

func TestButton_ToDomain(t *testing.T) {
	tests := []struct {
		name     string
		button   Button
		wantType d_message.ButtonType
	}{
		{
			name:     "postback button",
			button:   Button{Type: "postback", Title: "Click", Detail: "click_action"},
			wantType: d_message.POSTBACK,
		},
		{
			name:     "url button",
			button:   Button{Type: "url", Title: "Visit", Detail: "https://x.com"},
			wantType: d_message.URL,
		},
		{
			name:     "invalid type defaults to UNKNOWN",
			button:   Button{Type: "invalid", Title: "Test", Detail: "test"},
			wantType: d_message.UNKNOWN,
		},
		{
			name:     "empty type defaults to UNKNOWN",
			button:   Button{Type: "", Title: "Test", Detail: "test"},
			wantType: d_message.UNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.button.ToDomain()
			if got.Type != tt.wantType {
				t.Errorf("Button.ToDomain().Type = %v, want %v", got.Type, tt.wantType)
			}
			if got.Title != tt.button.Title {
				t.Errorf("Button.ToDomain().Title = %v, want %v", got.Title, tt.button.Title)
			}
			if got.Detail != tt.button.Detail {
				t.Errorf("Button.ToDomain().Detail = %v, want %v", got.Detail, tt.button.Detail)
			}
		})
	}
}

func TestMessage_ToDomain(t *testing.T) {
	tests := []struct {
		name    string
		message Message
	}{
		{
			name: "full message with buttons",
			message: Message{
				TextMessage: &TextMessage{
					ID:     "1",
					Title:  "Title",
					Detail: "Detail",
				},
				Buttons: []Button{
					{Type: "postback", Title: "Btn1", Detail: "btn1"},
					{Type: "url", Title: "Btn2", Detail: "https://x.com"},
				},
				DisplayButton: &Button{Type: "postback", Title: "Main", Detail: "main"},
				DateTime:      "2024-01-01T00:00:00Z",
				File:          &dto_file.File{ID: "f1", Type: "image", URL: "https://x.com/img.png"},
			},
		},
		{
			name: "message without buttons",
			message: Message{
				TextMessage:   &TextMessage{Detail: "Simple message"},
				DateTime:      "2024-01-01T00:00:00Z",
				DisplayButton: &Button{},
				File:          &dto_file.File{},
			},
		},
		{
			name: "empty message with nil pointers",
			message: Message{
				TextMessage:   &TextMessage{},
				DisplayButton: &Button{},
				File:          &dto_file.File{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.message.ToDomain()

			// Check buttons length
			if len(got.Buttons) != len(tt.message.Buttons) {
				t.Errorf("Message.ToDomain().Buttons length = %v, want %v", len(got.Buttons), len(tt.message.Buttons))
			}

			// Check DateTime
			if got.DateTime != tt.message.DateTime {
				t.Errorf("Message.ToDomain().DateTime = %v, want %v", got.DateTime, tt.message.DateTime)
			}
		})
	}
}
