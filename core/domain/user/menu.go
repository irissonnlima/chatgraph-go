package d_user

// Menu represents a menu configuration in the chatbot.
// Menus are associated with departments and can be active or inactive.
type Menu struct {
	// ID is the unique identifier for the menu.
	ID int
	// Name is the display name of the menu.
	Name string
}

// IsEmpty returns true if the Menu has an invalid ID (less than 1).
func (m Menu) IsEmpty() bool {
	return m.ID < 1
}
