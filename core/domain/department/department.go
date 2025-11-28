// Package d_department provides the Department domain model.
// Departments represent organizational units that can be hierarchically structured.
package d_department

// Department represents an organizational unit in the system.
// Departments can be nested using the ParentID field to create
// a hierarchical structure.
type Department struct {
	// ID is the unique identifier for the department.
	ID string
	// ParentID is the ID of the parent department, empty if this is a root department.
	ParentID string
	// Name is the display name of the department.
	Name string
	// Description provides additional details about the department.
	Description string
}
