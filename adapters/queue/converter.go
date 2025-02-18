package queue

type ChatIDJson struct {
	UserID    string `json:"user_id"`
	CompanyID string `json:"company_id"`
}

type UserStateJson struct {
	ChatID      ChatIDJson `json:"chat_id"`
	Menu        *string    `json:"menu"`
	Route       *string    `json:"route"`
	Observation *string    `json:"observation"`
	Protocol    *string    `json:"protocol"`
}

type MessageJson struct {
	UserState      UserStateJson `json:"user_state"`
	TypeMessage    string        `json:"type_message"`
	ContentMessage string        `json:"content_message"`
}
