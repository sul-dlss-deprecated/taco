package persistence

// Resource represents the resource as it exists in the persistence layer
type Resource struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
