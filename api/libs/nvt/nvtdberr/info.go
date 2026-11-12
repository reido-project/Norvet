package nvtdberr

// Info names the database objects an error is about. Empty fields mean the
// error message did not say. It never contains row values (e.g. the
// duplicated username), so it is safe to show to clients.
type Info struct {
	Database   string
	Table      string
	Column     string
	Constraint string // index / key / constraint name
}

// Subject returns the most specific name: column > table > database.
func (i Info) Subject() string {
	switch {
	case i.Column != "":
		return i.Column
	case i.Table != "":
		return i.Table
	default:
		return i.Database
	}
}
