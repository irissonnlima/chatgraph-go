// Package d_message provides message-related domain models for chat communication.
// It includes text messages, buttons, and file attachments.
package dto_message

import (
	dto_file "github.com/irissonnlima/chatgraph-go/adapters/dto/file"
	d_message "github.com/irissonnlima/chatgraph-go/core/domain/message"
)

// TextMessage represents the text content of a message.
type TextMessage struct {
	// ID is the unique identifier for the message.
	ID string `json:"id"`
	// Title is the main heading of the message.
	Title string `json:"title"`
	// Detail is the body text of the message.
	Detail string `json:"detail"`
	// Caption is an optional caption for media attachments.
	Caption string `json:"caption"`
	// MentionedIds contains IDs of users mentioned in the message.
	MentionedIds []string `json:"mentioned_ids"`
}

func (tm TextMessage) ToDomain() d_message.TextMessage {
	return d_message.TextMessage{
		ID:           tm.ID,
		Title:        tm.Title,
		Detail:       tm.Detail,
		Caption:      tm.Caption,
		MentionedIds: tm.MentionedIds,
	}
}

// Button represents an interactive button in a message.
type Button struct {
	// Type indicates whether this is a POSTBACK or URL button.
	Type string `json:"type"`
	// Title is the display text of the button.
	Title string `json:"title"`
	// Detail contains the postback value or URL depending on the type.
	Detail string `json:"detail"`
}

func (b Button) ToDomain() d_message.Button {
	t, err := d_message.ButtonTypeFromString(b.Type)
	if err != nil {
		t = d_message.UNKNOWN
	}

	return d_message.Button{
		Type:   t,
		Title:  b.Title,
		Detail: b.Detail,
	}
}

// Message represents a complete chat message with optional buttons and file attachments.
type Message struct {
	// TextMessage contains the text content of the message.
	TextMessage *TextMessage `json:"text_message"`
	// Buttons is a list of interactive buttons attached to the message.
	Buttons []Button `json:"buttons"`
	// DisplayButton is the primary action button displayed prominently.
	DisplayButton *Button `json:"display_button"`
	// DateTime is the timestamp when the message was sent or received.
	DateTime string `json:"date_time"`
	// File is an optional file attachment.
	File *dto_file.File `json:"file"`
}

func (m Message) ToDomain() d_message.Message {
	buttons := make([]d_message.Button, len(m.Buttons))
	for i, btn := range m.Buttons {
		buttons[i] = btn.ToDomain()
	}

	return d_message.Message{
		TextMessage:   m.TextMessage.ToDomain(),
		Buttons:       buttons,
		DisplayButton: m.DisplayButton.ToDomain(),
		DateTime:      m.DateTime,
		File:          m.File.ToDomain(),
	}
}
