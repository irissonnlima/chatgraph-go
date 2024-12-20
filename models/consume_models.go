package models

type JsonUserState struct {
	CustomerID  string `json:"customer_id"`
	Menu        string `json:"menu"`
	Route       string `json:"route"`
	Obs         string `json:"obs"`
	LstUpdate   string `json:"lst_update"`
	DirectionIn string `json:"direction_in"`
	VollId      string `json:"voll_id"`
	Platform    string `json:"platform"`
}

type JsonMessage struct {
	UserState     JsonUserState `json:"user_state"`
	CustomerPhone string        `json:"customer_phone"`
	CompanyPhone  string        `json:"company_phone"`
	Status        string        `json:"status"`
	Type          string        `json:"type"`
	Text          string        `json:"text"`
}
