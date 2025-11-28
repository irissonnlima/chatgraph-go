// Package d_message provides message-related domain models for chat communication.
// It includes text messages, buttons, and file attachments.
package d_message

import (
	d_file "chatgraph/core/domain/file"
	"fmt"
	"strings"
	"time"
)

// ButtonType represents the type of button in a message.
type ButtonType int

// String returns the string representation of the ButtonType.
func (bt ButtonType) String() string {
	switch bt {
	case POSTBACK:
		return "postback"
	case URL:
		return "url"
	default:
		return "unknown"
	}
}

// ButtonTypeFromString converts a string to a ButtonType.
// Returns an error if the string does not match a valid ButtonType.
func ButtonTypeFromString(s string) (ButtonType, error) {
	switch s {
	case "postback":
		return POSTBACK, nil
	case "url":
		return URL, nil
	default:
		return -1, fmt.Errorf("invalid button type: %s", s)
	}
}

// Maximum lengths for button fields.
const (
	// MAX_BUTTON_TITLE is the maximum length for a button title.
	MAX_BUTTON_TITLE = 20
	// MAX_BUTTON_DETAIL is the maximum length for a button detail.
	MAX_BUTTON_DETAIL = 30
)

// Button type constants.
const (
	// POSTBACK represents a button that sends a postback to the server.
	POSTBACK ButtonType = iota
	// URL represents a button that opens a URL.
	URL
)

// Error variables for button validation.
var (
	// ErrorButtonTypeInvalid is returned when the button type is not POSTBACK or URL.
	ErrorButtonTypeInvalid = fmt.Errorf("button type is invalid, must be either POSTBACK or URL")
	// ErrorButtonTitleTooLong is returned when the button title exceeds MAX_BUTTON_TITLE.
	ErrorButtonTitleTooLong = fmt.Errorf("button title is too long, maximum is %d characters", MAX_BUTTON_TITLE)
	// ErrorButtonDetailTooLong is returned when the button detail exceeds MAX_BUTTON_DETAIL.
	ErrorButtonDetailTooLong = fmt.Errorf("button detail is too long, maximum is %d characters", MAX_BUTTON_DETAIL)
)

// TextMessage represents the text content of a message.
type TextMessage struct {
	// ID is the unique identifier for the message.
	ID string
	// Title is the main heading of the message.
	Title string
	// Detail is the body text of the message.
	Detail string
	// Caption is an optional caption for media attachments.
	Caption string
	// MentionedIds contains IDs of users mentioned in the message.
	MentionedIds []string
}

// Button represents an interactive button in a message.
type Button struct {
	// Type indicates whether this is a POSTBACK or URL button.
	Type ButtonType
	// Title is the display text of the button.
	Title string
	// Detail contains the postback value or URL depending on the type.
	Detail string
}

// String returns a formatted text representation of the button.
// URL buttons display as "*Title*: URL" and POSTBACK buttons as "*Title*" or "_Detail_".
func (b Button) String() string {
	if b.Type == URL {
		return "\n*" + b.Title + "*: " + b.Detail
	}

	if b.Type == POSTBACK {
		if b.Title != "" {
			return "\n*" + b.Title + "*"
		}
		if b.Detail != "" {
			return "\n_" + b.Detail + "_"
		}
	}
	return ""
}

// Message represents a complete chat message with optional buttons and file attachments.
type Message struct {
	// TextMessage contains the text content of the message.
	TextMessage TextMessage
	// Buttons is a list of interactive buttons attached to the message.
	Buttons []Button
	// DisplayButton is the primary action button displayed prominently.
	DisplayButton Button
	// DateTime is the timestamp when the message was sent or received.
	DateTime time.Time
	// File is an optional file attachment.
	File d_file.File
}

// EntireText returns the complete text content of the message.
// It concatenates Title, Detail, button text representations, and Caption.
func (m Message) EntireText() string {
	buttons := m.TextButtons()

	return strings.Join(
		[]string{
			m.TextMessage.Title,
			m.TextMessage.Detail,
			buttons,
			m.TextMessage.Caption,
		},
		"\n",
	)
}

// HasButtons returns true if the message contains any buttons.
func (m Message) HasButtons() bool {
	return len(m.Buttons) > 0
}

// HasFile returns true if the message contains a file attachment.
func (m Message) HasFile() bool {
	return m.File.IsEmpty()
}

// ValidadeButtons validates all buttons in the message.
// Returns an error if any button has an invalid type or exceeds length limits.
func (m Message) ValidadeButtons() error {
	for _, button := range m.Buttons {
		if button.Type != POSTBACK && button.Type != URL {
			return ErrorButtonTypeInvalid
		}
		if len(button.Title) > MAX_BUTTON_TITLE {
			return ErrorButtonTitleTooLong
		}
		if len(button.Detail) > MAX_BUTTON_DETAIL {
			return ErrorButtonDetailTooLong
		}
	}
	return nil
}

// ClearAllButtons converts all buttons to text format and clears the button list.
// Both URL and POSTBACK buttons are converted to their text representation
// and appended to the message detail.
func (m *Message) TextButtons() string {
	result := ""
	for _, b := range m.Buttons {
		result += strings.Join([]string{result, b.String()}, "\n")
	}
	return result
}
