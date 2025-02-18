package domain_primitives

import (
	"strings"
)

func NilString(stringPointer *string) string {
	if stringPointer == nil {
		return ""
	}

	return *stringPointer
}

type ChatID struct {
	UserID    string
	CompanyID string
}

func (c ChatID) Stringfy() string {
	limiter := strings.Repeat("-", 30) + "\n"
	return limiter + "UserID: " + c.UserID + ", CompanyID: " + c.CompanyID + "\n"
}

type UserState struct {
	ChatID      ChatID
	Route       Router
	Menu        string
	Observation *string
	Protocol    *string
}

func (u UserState) Stringfy() string {
	chatIDString := u.ChatID.Stringfy()
	uStateString := chatIDString + "Menu: " + u.Menu + "\n" +
		"Route: " + u.Route.HistoryRoute() + "\n" +
		"Observation: " + NilString(u.Observation) + "\n" +
		"Protocol: " + NilString(u.Protocol) + "\n"

	return uStateString
}

type UserCall struct {
	UserState UserState
	Message   Message
}

func (u UserCall) Stringfy() string {
	userStateString := u.UserState.Stringfy()
	messageString := u.Message.Stringfy()
	return userStateString + messageString
}
