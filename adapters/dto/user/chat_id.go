package dto_user

import d_user "github.com/irissonnlima/chatgraph-go/core/domain/user"

// ChatID uniquely identifies a chat session by combining user and company identifiers.
// Both fields are required for a valid ChatID.
type ChatID struct {
	// UserID is the unique identifier for the user.
	UserID string `json:"user_id"`
	// CompanyID is the unique identifier for the company.
	CompanyID string `json:"company_id"`
}

func (c ChatID) ToDomain() d_user.ChatID {
	return d_user.ChatID{
		UserID:    c.UserID,
		CompanyID: c.CompanyID,
	}
}
