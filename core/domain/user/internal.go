package d_user

// UserInternal represents internal HR data for the user.
// This is populated only when the user has an internal employment relationship.
type UserInternal struct {
	// Matricula is the employee registration number.
	Matricula string
	// Cargo is the employee's job title.
	Cargo string
	// Filial is the employee's branch.
	Filial string
	// Empresa is the employee's company name.
	Empresa string
	// DataAdmissao is the employee's admission date.
	DataAdmissao string
}
