package d_user

// UserIdentity represents the authenticated identity of the user.
type UserIdentity struct {
	// AuthLevel is the user's authorization level.
	AuthLevel AuthLevel
	// CPF is the authenticated CPF (Brazilian taxpayer ID).
	CPF string
	// Active indicates whether the user is active.
	Active bool
	// AuthStatus is the current authentication status string.
	AuthStatus string
	// DeviceID is the identifier of the user's device.
	DeviceID string
}
