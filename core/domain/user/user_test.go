package d_user

import (
	d_route "chatgraph/core/domain/route"
	"testing"
)

// Test structs for generic UserState tests
type TestObservation struct {
	OrderID   string `json:"order_id"`
	ProductID int    `json:"product_id"`
	Active    bool   `json:"active"`
}

type EmptyObservation struct{}

// ============================================
// ChatID Tests
// ============================================

func TestChatID_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		chatID   ChatID
		expected bool
	}{
		{
			name:     "empty ChatID - both fields empty",
			chatID:   ChatID{},
			expected: true,
		},
		{
			name:     "empty UserID only",
			chatID:   ChatID{UserID: "", CompanyID: "company-123"},
			expected: true,
		},
		{
			name:     "empty CompanyID only",
			chatID:   ChatID{UserID: "user-456", CompanyID: ""},
			expected: true,
		},
		{
			name:     "both fields filled",
			chatID:   ChatID{UserID: "user-456", CompanyID: "company-123"},
			expected: false,
		},
		{
			name:     "whitespace UserID",
			chatID:   ChatID{UserID: "   ", CompanyID: "company-123"},
			expected: false, // whitespace is not empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.chatID.IsEmpty()
			if got != tt.expected {
				t.Errorf("ChatID.IsEmpty() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// ============================================
// Menu Tests
// ============================================

func TestMenu_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		menu     Menu
		expected bool
	}{
		{
			name:     "empty Menu - zero value",
			menu:     Menu{},
			expected: true,
		},
		{
			name:     "ID is zero",
			menu:     Menu{ID: 0, Name: "Main Menu"},
			expected: true,
		},
		{
			name:     "ID is negative",
			menu:     Menu{ID: -1, Name: "Invalid Menu"},
			expected: true,
		},
		{
			name:     "ID is 1 - valid",
			menu:     Menu{ID: 1, Name: "Main Menu"},
			expected: false,
		},
		{
			name:     "ID greater than 1",
			menu:     Menu{ID: 100, Name: "Settings"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.menu.IsEmpty()
			if got != tt.expected {
				t.Errorf("Menu.IsEmpty() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// ============================================
// User Tests
// ============================================

func TestUser_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		user     User
		expected bool
	}{
		{
			name:     "empty User - zero value",
			user:     User{},
			expected: true,
		},
		{
			name:     "only CPF filled",
			user:     User{CPF: "12345678900"},
			expected: false,
		},
		{
			name:     "only Name filled",
			user:     User{Name: "John Doe"},
			expected: false,
		},
		{
			name:     "only Phone filled",
			user:     User{Phone: "+5511999999999"},
			expected: false,
		},
		{
			name:     "only Email filled",
			user:     User{Email: "john@example.com"},
			expected: false,
		},
		{
			name:     "all fields filled",
			user:     User{CPF: "12345678900", Name: "John Doe", Phone: "+5511999999999", Email: "john@example.com"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.IsEmpty()
			if got != tt.expected {
				t.Errorf("User.IsEmpty() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// ============================================
// UserState Tests
// ============================================

func TestUserState_IsEmpty(t *testing.T) {
	tests := []struct {
		name      string
		userState UserState[TestObservation]
		expected  bool
	}{
		{
			name:      "empty UserState - zero value",
			userState: UserState[TestObservation]{},
			expected:  true,
		},
		{
			name: "ChatID with empty UserID",
			userState: UserState[TestObservation]{
				ChatID: ChatID{UserID: "", CompanyID: "company-123"},
			},
			expected: true,
		},
		{
			name: "ChatID with empty CompanyID",
			userState: UserState[TestObservation]{
				ChatID: ChatID{UserID: "user-456", CompanyID: ""},
			},
			expected: true,
		},
		{
			name: "valid ChatID - not empty",
			userState: UserState[TestObservation]{
				ChatID: ChatID{UserID: "user-456", CompanyID: "company-123"},
			},
			expected: false,
		},
		{
			name: "full UserState with valid ChatID",
			userState: UserState[TestObservation]{
				SessionID: 12345,
				ChatID:    ChatID{UserID: "user-456", CompanyID: "company-123"},
				User: User{
					CPF:   "12345678900",
					Name:  "John Doe",
					Phone: "+5511999999999",
					Email: "john@example.com",
				},
				Menu: Menu{
					ID:   1,
					Name: "Main Menu",
				},
				Route:       d_route.NewRoute("start.menu", '.'),
				DirectionIn: true,
				Observation: TestObservation{OrderID: "ORD-123", ProductID: 456, Active: true},
				Platform:    "whatsapp",
				DtCreated:   "2025-11-27T09:00:00Z",
			},
			expected: false,
		},
		{
			name: "UserState with other fields but empty ChatID",
			userState: UserState[TestObservation]{
				SessionID:   99999,
				User:        User{Name: "Jane Doe"},
				Menu:        Menu{ID: 5},
				DirectionIn: true,
				Platform:    "telegram",
			},
			expected: true, // IsEmpty depends only on ChatID
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.userState.IsEmpty()
			if got != tt.expected {
				t.Errorf("UserState.IsEmpty() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestUserState_IsEmpty_WithDifferentTypes(t *testing.T) {
	t.Run("UserState with EmptyObservation", func(t *testing.T) {
		state := UserState[EmptyObservation]{
			ChatID: ChatID{UserID: "user-123", CompanyID: "company-456"},
		}
		if state.IsEmpty() {
			t.Error("expected UserState to not be empty")
		}
	})

	t.Run("UserState with string observation", func(t *testing.T) {
		state := UserState[string]{
			ChatID:      ChatID{UserID: "user-123", CompanyID: "company-456"},
			Observation: "simple string observation",
		}
		if state.IsEmpty() {
			t.Error("expected UserState to not be empty")
		}
	})

	t.Run("UserState with map observation", func(t *testing.T) {
		state := UserState[map[string]any]{
			ChatID:      ChatID{UserID: "user-123", CompanyID: "company-456"},
			Observation: map[string]any{"key": "value"},
		}
		if state.IsEmpty() {
			t.Error("expected UserState to not be empty")
		}
	})
}

func TestUserState_LoadObservation(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expected    TestObservation
		expectError bool
	}{
		{
			name: "valid JSON with all fields",
			json: `{"order_id": "ORD-001", "product_id": 123, "active": true}`,
			expected: TestObservation{
				OrderID:   "ORD-001",
				ProductID: 123,
				Active:    true,
			},
			expectError: false,
		},
		{
			name: "valid JSON with partial fields",
			json: `{"order_id": "ORD-002"}`,
			expected: TestObservation{
				OrderID:   "ORD-002",
				ProductID: 0,
				Active:    false,
			},
			expectError: false,
		},
		{
			name:        "empty JSON object",
			json:        `{}`,
			expected:    TestObservation{},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			json:        `{invalid json}`,
			expected:    TestObservation{},
			expectError: true,
		},
		{
			name:        "empty string",
			json:        ``,
			expected:    TestObservation{},
			expectError: true,
		},
		{
			name: "JSON with extra fields (ignored)",
			json: `{"order_id": "ORD-003", "product_id": 999, "active": false, "extra": "ignored"}`,
			expected: TestObservation{
				OrderID:   "ORD-003",
				ProductID: 999,
				Active:    false,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &UserState[TestObservation]{}
			err := state.LoadObservation(tt.json)

			if (err != nil) != tt.expectError {
				t.Errorf("LoadObservation() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError {
				if state.Observation != tt.expected {
					t.Errorf("LoadObservation() observation = %+v, expected %+v", state.Observation, tt.expected)
				}
			}
		})
	}
}

func TestUserState_LoadObservation_PreservesOtherFields(t *testing.T) {
	state := &UserState[TestObservation]{
		SessionID: 12345,
		ChatID:    ChatID{UserID: "user-456", CompanyID: "company-123"},
		Platform:  "whatsapp",
	}

	err := state.LoadObservation(`{"order_id": "ORD-999", "product_id": 42, "active": true}`)
	if err != nil {
		t.Fatalf("LoadObservation() unexpected error: %v", err)
	}

	// Verify observation was loaded
	expectedObs := TestObservation{OrderID: "ORD-999", ProductID: 42, Active: true}
	if state.Observation != expectedObs {
		t.Errorf("Observation = %+v, expected %+v", state.Observation, expectedObs)
	}

	// Verify other fields were preserved
	if state.SessionID != 12345 {
		t.Errorf("SessionID = %d, expected 12345", state.SessionID)
	}
	if state.ChatID.UserID != "user-456" {
		t.Errorf("ChatID.UserID = %s, expected user-456", state.ChatID.UserID)
	}
	if state.Platform != "whatsapp" {
		t.Errorf("Platform = %s, expected whatsapp", state.Platform)
	}
}
