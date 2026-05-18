// Package d_user provides user-related domain models including User, UserState, ChatID, and Menu.
// These models represent the user's identity, session state, and navigation context.
package dto_user

import (
	"encoding/json"
	"log"

	d_route "github.com/irissonnlima/chatgraph-go/core/domain/route"
	d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"
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
	// Identity holds the user's authentication identity data.
	Identity *UserIdentity `json:"identity,omitempty"`
	// Internal holds internal HR data. Nil if user has no internal relationship.
	Internal *UserInternal `json:"internal,omitempty"`
}

// UserIdentity represents the authenticated identity of the user (DTO).
type UserIdentity struct {
	AuthLevel  string `json:"auth_level"`
	CPF        string `json:"cpf"`
	Active     bool   `json:"active"`
	AuthStatus string `json:"auth_status"`
	DeviceID   string `json:"device_id"`
}

// UserInternal represents internal HR data (DTO).
type UserInternal struct {
	Matricula    string `json:"matricula"`
	Cargo        string `json:"cargo"`
	Filial       string `json:"filial"`
	Empresa      string `json:"empresa"`
	DataAdmissao string `json:"data_admissao"`
}

func (u User) ToDomain() d_user.User {
	user := d_user.User{
		CPF:               u.CPF,
		AuthorizationCode: u.AuthorizationCode,
		Name:              u.Name,
		Phone:             u.Phone,
		Email:             u.Email,
	}

	if u.Identity != nil {
		user.Identity = d_user.UserIdentity{
			AuthLevel:  d_user.ParseAuthLevel(u.Identity.AuthLevel),
			CPF:        u.Identity.CPF,
			Active:     u.Identity.Active,
			AuthStatus: u.Identity.AuthStatus,
			DeviceID:   u.Identity.DeviceID,
		}
	}

	if u.Internal != nil {
		user.Internal = &d_user.UserInternal{
			Matricula:    u.Internal.Matricula,
			Cargo:        u.Internal.Cargo,
			Filial:       u.Internal.Filial,
			Empresa:      u.Internal.Empresa,
			DataAdmissao: u.Internal.DataAdmissao,
		}
	}

	return user
}

// UserState represents the complete state of a user's chat session.
// It is generic over Obs, which allows custom observation data to be
// associated with the session.
type UserState struct {
	// SessionID is the unique identifier for the current session.
	SessionID int64 `json:"session_id,omitempty"`
	// ChatID identifies the user and company for this chat.
	ChatID *ChatID `json:"chat_id,omitempty"`
	// User contains the user's personal information.
	User *User `json:"user,omitempty"`
	// Menu is the current menu context.
	Menu *Menu `json:"menu,omitempty"`
	// Route tracks the navigation history through the chatbot.
	Route string `json:"route,omitempty"`
	// DirectionIn indicates if the message is incoming (true) or outgoing (false).
	DirectionIn bool `json:"direction_in,omitempty"`
	// Observation holds custom session data of type Obs.
	Observation string `json:"observation,omitempty"`
	// Platform identifies the messaging platform (e.g., "whatsapp", "telegram").
	Platform string `json:"platform,omitempty"`
	// LastUpdate is the timestamp of the last state update.
	LastUpdate string `json:"last_update,omitempty"`
	// DtCreated is the timestamp when the session was created.
	DtCreated string `json:"dt_created,omitempty"`
}

func UserStateToDomain[Obs any](u UserState) d_user.UserState[Obs] {
	var obs Obs
	if u.Observation != "" {
		err := json.Unmarshal([]byte(u.Observation), &obs)
		if err != nil {
			log.Println(err)
			panic("failed to unmarshal observation")
		}
	}

	state := d_user.UserState[Obs]{
		SessionID:   u.SessionID,
		Route:       d_route.NewRoute(u.Route, '.'),
		DirectionIn: u.DirectionIn,
		Observation: obs,
		Platform:    u.Platform,
		DtCreated:   u.DtCreated,
	}

	if u.ChatID != nil {
		state.ChatID = u.ChatID.ToDomain()
	}
	if u.User != nil {
		state.User = u.User.ToDomain()
	}
	if u.Menu != nil {
		state.Menu = u.Menu.ToDomain()
	}

	return state
}
