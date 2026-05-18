package d_user

import "testing"

func TestAuthLevel_String(t *testing.T) {
	tests := []struct {
		level    AuthLevel
		expected string
	}{
		{AuthBlocked, "blocked"},
		{AuthUnknown, "unknown"},
		{AuthRead, "read"},
		{AuthWrite, "write"},
		{AuthLevel(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("AuthLevel(%d).String() = %v, want %v", tt.level, got, tt.expected)
			}
		})
	}
}

func TestParseAuthLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected AuthLevel
	}{
		{"blocked", AuthBlocked},
		{"unknown", AuthUnknown},
		{"read", AuthRead},
		{"write", AuthWrite},
		{"invalid", AuthUnknown},
		{"", AuthUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseAuthLevel(tt.input); got != tt.expected {
				t.Errorf("ParseAuthLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestAuthLevel_IsAtLeast(t *testing.T) {
	tests := []struct {
		name     string
		level    AuthLevel
		required AuthLevel
		expected bool
	}{
		{"write >= write", AuthWrite, AuthWrite, true},
		{"write >= read", AuthWrite, AuthRead, true},
		{"write >= unknown", AuthWrite, AuthUnknown, true},
		{"write >= blocked", AuthWrite, AuthBlocked, true},
		{"read >= read", AuthRead, AuthRead, true},
		{"read >= write", AuthRead, AuthWrite, false},
		{"unknown >= read", AuthUnknown, AuthRead, false},
		{"blocked >= unknown", AuthBlocked, AuthUnknown, false},
		{"blocked >= blocked", AuthBlocked, AuthBlocked, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.IsAtLeast(tt.required); got != tt.expected {
				t.Errorf("AuthLevel(%d).IsAtLeast(%d) = %v, want %v", tt.level, tt.required, got, tt.expected)
			}
		})
	}
}

func TestAuthLevel_Ordering(t *testing.T) {
	if AuthBlocked >= AuthUnknown {
		t.Error("AuthBlocked should be less than AuthUnknown")
	}
	if AuthUnknown >= AuthRead {
		t.Error("AuthUnknown should be less than AuthRead")
	}
	if AuthRead >= AuthWrite {
		t.Error("AuthRead should be less than AuthWrite")
	}
}
