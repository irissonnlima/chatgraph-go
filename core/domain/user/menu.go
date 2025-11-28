package d_user

// Menu represents a menu configuration in the chatbot.
// Menus are associated with departments and can be active or inactive.
type Menu struct {
	// ID is the unique identifier for the menu.
	ID int
	// DepartmentID is the department this menu belongs to.
	DepartmentID int
	// Name is the display name of the menu.
	Name string
	// Description provides additional details about the menu.
	Description string
	// Active indicates whether the menu is currently available.
	Active bool
}

// IsEmpty returns true if the Menu has an invalid ID (less than 1).
func (m Menu) IsEmpty() bool {
	return m.ID < 1
}
