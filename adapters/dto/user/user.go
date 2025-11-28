// Package d_user provides user-related domain models including User, UserState, ChatID, and Menu.
// These models represent the user's identity, session state, and navigation context.
package dto_user

import (
	d_route "chatgraph/core/domain/route"
	d_user "chatgraph/core/domain/user"
	"encoding/json"
	"log"
)

// User represents the basic user information.
// This contains the user's personal data that persists across sessions.
type User struct {
	// CPF is the Brazilian individual taxpayer registry identification.
	CPF string `json:"cpf"`
	// AuthorizationCode is a code used for user authorization. This indicates if the user is authenticated.
	AuthorizationCode string `json:"authorization_code"`
	// Name is the user's full name.
	Name string `json:"name"`
	// Phone is the user's phone number.
	Phone string `json:"phone"`
	// Email is the user's email address.
	Email string `json:"email"`
}

func (u User) ToDomain() d_user.User {
	return d_user.User{
		CPF:               u.CPF,
		AuthorizationCode: u.AuthorizationCode,
		Name:              u.Name,
		Phone:             u.Phone,
		Email:             u.Email,
	}
}

// UserState represents the complete state of a user's chat session.
// It is generic over Obs, which allows custom observation data to be
// associated with the session.
type UserState struct {
	// SessionID is the unique identifier for the current session.
	SessionID int64 `json:"session_id"`
	// ChatID identifies the user and company for this chat.
	ChatID ChatID `json:"chat_id"`
	// User contains the user's personal information.
	User User `json:"user"`
	// Menu is the current menu context.
	Menu Menu `json:"menu"`
	// Route tracks the navigation history through the chatbot.
	Route string `json:"route"`
	// DirectionIn indicates if the message is incoming (true) or outgoing (false).
	DirectionIn bool `json:"direction_in"`
	// Observation holds custom session data of type Obs.
	Observation string `json:"observation"`
	// Platform identifies the messaging platform (e.g., "whatsapp", "telegram").
	Platform string `json:"platform"`
	// LastUpdate is the timestamp of the last state update.
	LastUpdate string `json:"last_update"`
	// DtCreated is the timestamp when the session was created.
	DtCreated string `json:"dt_created"`
}

func UserStateToDomain[Obs any](u UserState) d_user.UserState[Obs] {
	var obs Obs
	err := json.Unmarshal([]byte(u.Observation), &obs)
	if err != nil {
		log.Println(err)
		panic("failed to unmarshal observation")
	}

	return d_user.UserState[Obs]{
		SessionID:   u.SessionID,
		ChatID:      u.ChatID.ToDomain(),
		User:        u.User.ToDomain(),
		Menu:        u.Menu.ToDomain(),
		Route:       d_route.NewRoute(u.Route, '.'),
		DirectionIn: u.DirectionIn,
		Observation: obs,
		Platform:    u.Platform,
		DtCreated:   u.DtCreated,
	}
}
