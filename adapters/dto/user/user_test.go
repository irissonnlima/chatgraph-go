package dto_user

import (
	"testing"
)

func TestChatID_ToDomain(t *testing.T) {
	tests := []struct {
		name   string
		chatID ChatID
	}{
		{
			name:   "valid chat id",
			chatID: ChatID{UserID: "user123", CompanyID: "company456"},
		},
		{
			name:   "empty chat id",
			chatID: ChatID{},
		},
		{
			name:   "partial chat id",
			chatID: ChatID{UserID: "user123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.chatID.ToDomain()
			if got.UserID != tt.chatID.UserID {
				t.Errorf("ChatID.ToDomain().UserID = %v, want %v", got.UserID, tt.chatID.UserID)
			}
			if got.CompanyID != tt.chatID.CompanyID {
				t.Errorf("ChatID.ToDomain().CompanyID = %v, want %v", got.CompanyID, tt.chatID.CompanyID)
			}
		})
	}
}

func TestUser_ToDomain(t *testing.T) {
	tests := []struct {
		name string
		user User
	}{
		{
			name: "full user",
			user: User{
				CPF:               "12345678900",
				AuthorizationCode: "auth123",
				Name:              "John Doe",
				Phone:             "11999999999",
				Email:             "john@example.com",
			},
		},
		{
			name: "empty user",
			user: User{},
		},
		{
			name: "partial user",
			user: User{
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.ToDomain()
			if got.CPF != tt.user.CPF {
				t.Errorf("User.ToDomain().CPF = %v, want %v", got.CPF, tt.user.CPF)
			}
			if got.AuthorizationCode != tt.user.AuthorizationCode {
				t.Errorf("User.ToDomain().AuthorizationCode = %v, want %v", got.AuthorizationCode, tt.user.AuthorizationCode)
			}
			if got.Name != tt.user.Name {
				t.Errorf("User.ToDomain().Name = %v, want %v", got.Name, tt.user.Name)
			}
			if got.Phone != tt.user.Phone {
				t.Errorf("User.ToDomain().Phone = %v, want %v", got.Phone, tt.user.Phone)
			}
			if got.Email != tt.user.Email {
				t.Errorf("User.ToDomain().Email = %v, want %v", got.Email, tt.user.Email)
			}
		})
	}
}

func TestMenu_ToDomain(t *testing.T) {
	tests := []struct {
		name string
		menu Menu
	}{
		{
			name: "full menu",
			menu: Menu{ID: 1, Name: "Main Menu", Description: "Main menu description"},
		},
		{
			name: "empty menu",
			menu: Menu{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.menu.ToDomain()
			if got.ID != tt.menu.ID {
				t.Errorf("Menu.ToDomain().ID = %v, want %v", got.ID, tt.menu.ID)
			}
			if got.Name != tt.menu.Name {
				t.Errorf("Menu.ToDomain().Name = %v, want %v", got.Name, tt.menu.Name)
			}
		})
	}
}

func TestUserStateToDomain(t *testing.T) {
	type TestObs struct {
		Value string `json:"value"`
	}

	tests := []struct {
		name      string
		userState UserState
		wantPanic bool
	}{
		{
			name: "full user state",
			userState: UserState{
				SessionID:   123,
				ChatID:      &ChatID{UserID: "user1", CompanyID: "comp1"},
				User:        &User{Name: "John"},
				Menu:        &Menu{ID: 1, Name: "Main"},
				Route:       "menu.submenu",
				DirectionIn: true,
				Observation: `{"value":"test"}`,
				Platform:    "whatsapp",
				DtCreated:   "2024-01-01",
			},
			wantPanic: false,
		},
		{
			name: "minimal user state",
			userState: UserState{
				SessionID: 456,
				Route:     "start",
			},
			wantPanic: false,
		},
		{
			name: "user state with empty observation",
			userState: UserState{
				SessionID: 789,
			},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("UserStateToDomain() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()

			got := UserStateToDomain[TestObs](tt.userState)

			if got.SessionID != tt.userState.SessionID {
				t.Errorf("UserStateToDomain().SessionID = %v, want %v", got.SessionID, tt.userState.SessionID)
			}
			if got.DirectionIn != tt.userState.DirectionIn {
				t.Errorf("UserStateToDomain().DirectionIn = %v, want %v", got.DirectionIn, tt.userState.DirectionIn)
			}
			if got.Platform != tt.userState.Platform {
				t.Errorf("UserStateToDomain().Platform = %v, want %v", got.Platform, tt.userState.Platform)
			}
			if got.DtCreated != tt.userState.DtCreated {
				t.Errorf("UserStateToDomain().DtCreated = %v, want %v", got.DtCreated, tt.userState.DtCreated)
			}
		})
	}
}

func TestUserStateToDomain_InvalidObservation(t *testing.T) {
	type TestObs struct {
		Value string `json:"value"`
	}

	userState := UserState{
		SessionID:   123,
		Observation: "invalid json",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("UserStateToDomain() should panic with invalid observation JSON")
		}
	}()

	UserStateToDomain[TestObs](userState)
}
