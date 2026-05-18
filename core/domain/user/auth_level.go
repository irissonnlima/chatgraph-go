package d_user

// AuthLevel represents the user's authorization level.
// Levels are ordered: Blocked < Unknown < Read < Write.
type AuthLevel int

const (
	// AuthBlocked indicates the user is blocked from accessing the system.
	AuthBlocked AuthLevel = iota
	// AuthUnknown indicates the user's authorization level is not yet determined.
	AuthUnknown
	// AuthRead indicates the user has read-only access.
	AuthRead
	// AuthWrite indicates the user has full read-write access.
	AuthWrite
)

// String returns the string representation of the AuthLevel.
func (a AuthLevel) String() string {
	switch a {
	case AuthBlocked:
		return "blocked"
	case AuthUnknown:
		return "unknown"
	case AuthRead:
		return "read"
	case AuthWrite:
		return "write"
	default:
		return "unknown"
	}
}

// ParseAuthLevel converts a string to an AuthLevel.
// Returns AuthUnknown for unrecognized strings.
func ParseAuthLevel(s string) AuthLevel {
	switch s {
	case "blocked":
		return AuthBlocked
	case "unknown":
		return AuthUnknown
	case "read":
		return AuthRead
	case "write":
		return AuthWrite
	default:
		return AuthUnknown
	}
}

// IsAtLeast returns true if this level is greater than or equal to the required level.
func (a AuthLevel) IsAtLeast(required AuthLevel) bool {
	return a >= required
}
