package authorization

// Agent represents a user or machine account that uses the api
type Agent struct {
	// Identifier is typically an email address of the form: lmcrae@stanford.edu
	// but it could also include special identifiers, e.g.: labs@stanford
	Identifier string
}

func (d *Agent) String() string {
	return d.Identifier
}
