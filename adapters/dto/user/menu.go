package dto_user

import d_user "chatgraph/core/domain/user"

// Menu represents a menu configuration in the chatbot.
// Menus are associated with departments and can be active or inactive.
type Menu struct {
	// ID is the unique identifier for the menu.
	ID int `json:"id"`
	// Name is the display name of the menu.
	Name string `json:"name"`
	// Description provides additional details about the menu.
	Description string `json:"description"`
}

func (m Menu) ToDomain() d_user.Menu {
	return d_user.Menu{
		ID:   m.ID,
		Name: m.Name,
	}
}
