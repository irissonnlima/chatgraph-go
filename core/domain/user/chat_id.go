package d_user

// ChatID uniquely identifies a chat session by combining user and company identifiers.
// Both fields are required for a valid ChatID.
type ChatID struct {
	// UserID is the unique identifier for the user.
	UserID string
	// CompanyID is the unique identifier for the company.
	CompanyID string
}

// IsEmpty returns true if either UserID or CompanyID is empty.
// A valid ChatID must have both fields populated.
func (c ChatID) IsEmpty() bool {
	return c.CompanyID == "" || c.UserID == ""
}
