// Package d_user provides user-related domain models including User, UserState, ChatID, and Menu.
// These models represent the user's identity, session state, and navigation context.
package d_user

import (
	d_route "chatgraph/core/domain/route"
	"encoding/json"
)

// User represents the basic user information.
// This contains the user's personal data that persists across sessions.
type User struct {
	// CPF is the Brazilian individual taxpayer registry identification.
	CPF string
	// AuthorizationCode is a code used for user authorization. This indicates if the user is authenticated.
	AuthorizationCode string
	// Name is the user's full name.
	Name string
	// Phone is the user's phone number.
	Phone string
	// Email is the user's email address.
	Email string
}

// IsEmpty returns true if the User is a zero-value struct.
func (u User) IsEmpty() bool {
	return u == (User{})
}

// UserState represents the complete state of a user's chat session.
// It is generic over Obs, which allows custom observation data to be
// associated with the session.
type UserState[Obs any] struct {
	// SessionID is the unique identifier for the current session.
	SessionID int64
	// ChatID identifies the user and company for this chat.
	ChatID ChatID
	// User contains the user's personal information.
	User User
	// Menu is the current menu context.
	Menu Menu
	// Route tracks the navigation history through the chatbot.
	Route d_route.Route
	// DirectionIn indicates if the message is incoming (true) or outgoing (false).
	DirectionIn bool
	// Observation holds custom session data of type Obs.
	Observation Obs
	// Platform identifies the messaging platform (e.g., "whatsapp", "telegram").
	Platform string
	// LastUpdate is the timestamp of the last state update.
	LastUpdate string
	// DtCreated is the timestamp when the session was created.
	DtCreated string
}

// IsEmpty returns true if the UserState has an empty ChatID.
func (u UserState[Obs]) IsEmpty() bool {
	return u.ChatID.IsEmpty()
}

// LoadObservation deserializes a JSON string into the Observation field.
// Returns an error if the JSON is invalid or cannot be unmarshaled into type Obs.
func (u *UserState[Obs]) LoadObservation(observation string) error {
	var obs Obs

	err := json.Unmarshal([]byte(observation), &obs)
	if err != nil {
		return err
	}

	u.Observation = obs
	return nil
}
