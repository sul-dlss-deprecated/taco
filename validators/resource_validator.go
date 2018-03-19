package validators

// ResourceValidator is the interface for validators that check the resources format
type ResourceValidator interface {
	ValidateResource(body string) error
}
