package d_message

import (
	"testing"
)

func TestButtonType_String(t *testing.T) {
	tests := []struct {
		name       string
		buttonType ButtonType
		want       string
	}{
		{"POSTBACK returns postback", POSTBACK, "postback"},
		{"URL returns url", URL, "url"},
		{"UNKNOWN returns unknown", UNKNOWN, "unknown"},
		{"invalid returns unknown", ButtonType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.buttonType.String(); got != tt.want {
				t.Errorf("ButtonType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButtonTypeFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ButtonType
		wantErr bool
	}{
		{"postback string", "postback", POSTBACK, false},
		{"url string", "url", URL, false},
		{"invalid string", "invalid", -1, true},
		{"empty string", "", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ButtonTypeFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ButtonTypeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ButtonTypeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButton_String(t *testing.T) {
	tests := []struct {
		name   string
		button Button
		want   string
	}{
		{
			name:   "URL button with title",
			button: Button{Type: URL, Title: "Visit", Detail: "https://example.com"},
			want:   "\n*Visit*: https://example.com",
		},
		{
			name:   "POSTBACK button with title",
			button: Button{Type: POSTBACK, Title: "Click Me", Detail: "click"},
			want:   "\n*Click Me*",
		},
		{
			name:   "POSTBACK button without title but with detail",
			button: Button{Type: POSTBACK, Title: "", Detail: "some_detail"},
			want:   "\n_some_detail_",
		},
		{
			name:   "POSTBACK button without title and detail",
			button: Button{Type: POSTBACK, Title: "", Detail: ""},
			want:   "",
		},
		{
			name:   "UNKNOWN button",
			button: Button{Type: UNKNOWN, Title: "Test", Detail: "test"},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.button.String(); got != tt.want {
				t.Errorf("Button.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestButton_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		button Button
		want   bool
	}{
		{"empty title", Button{Type: POSTBACK, Title: "", Detail: "test"}, true},
		{"UNKNOWN type", Button{Type: UNKNOWN, Title: "Test", Detail: "test"}, true},
		{"valid button", Button{Type: POSTBACK, Title: "Test", Detail: "test"}, false},
		{"URL button valid", Button{Type: URL, Title: "Link", Detail: "https://x.com"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.button.IsEmpty(); got != tt.want {
				t.Errorf("Button.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_EntireText(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		want    string
	}{
		{
			name: "all fields filled",
			message: Message{
				TextMessage: TextMessage{
					Title:   "Title",
					Detail:  "Detail",
					Caption: "Caption",
				},
			},
			want: "Title\nDetail\nCaption",
		},
		{
			name: "only detail",
			message: Message{
				TextMessage: TextMessage{
					Detail: "Just detail",
				},
			},
			want: "Just detail",
		},
		{
			name: "only title and detail",
			message: Message{
				TextMessage: TextMessage{
					Title:  "Title",
					Detail: "Detail",
				},
			},
			want: "Title\nDetail",
		},
		{
			name:    "empty message",
			message: Message{},
			want:    "",
		},
		{
			name: "title only",
			message: Message{
				TextMessage: TextMessage{
					Title: "Only Title",
				},
			},
			want: "Only Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.message.EntireText(); got != tt.want {
				t.Errorf("Message.EntireText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMessage_HasButtons(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		want    bool
	}{
		{"no buttons", Message{}, false},
		{"with buttons", Message{Buttons: []Button{{Type: POSTBACK, Title: "btn"}}}, true},
		{"empty buttons slice", Message{Buttons: []Button{}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.message.HasButtons(); got != tt.want {
				t.Errorf("Message.HasButtons() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_HasFile(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		want    bool
	}{
		{"no file (empty)", Message{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.message.HasFile(); got != tt.want {
				t.Errorf("Message.HasFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_ValidadeButtons(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		wantErr bool
	}{
		{"no buttons", Message{}, false},
		{"valid POSTBACK buttons", Message{Buttons: []Button{{Type: POSTBACK, Title: "btn"}}}, false},
		{"valid URL buttons", Message{Buttons: []Button{{Type: URL, Title: "link", Detail: "https://x.com"}}}, false},
		{"invalid button type", Message{Buttons: []Button{{Type: UNKNOWN, Title: "btn"}}}, true},
		{"mixed valid buttons", Message{Buttons: []Button{
			{Type: POSTBACK, Title: "btn"},
			{Type: URL, Title: "link"},
		}}, false},
		{"one invalid among valid", Message{Buttons: []Button{
			{Type: POSTBACK, Title: "btn"},
			{Type: UNKNOWN, Title: "invalid"},
		}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.ValidadeButtons()
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.ValidadeButtons() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessage_TextButtons(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		want    string
	}{
		{"no buttons", Message{}, ""},
		{
			"single POSTBACK button",
			Message{Buttons: []Button{{Type: POSTBACK, Title: "Click"}}},
			"\n\n*Click*",
		},
		{
			"single URL button",
			Message{Buttons: []Button{{Type: URL, Title: "Visit", Detail: "https://x.com"}}},
			"\n\n*Visit*: https://x.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.message.TextButtons()
			if got != tt.want {
				t.Errorf("Message.TextButtons() = %q, want %q", got, tt.want)
			}
		})
	}
}
