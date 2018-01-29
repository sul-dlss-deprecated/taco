package persistence

// Resource represents the resource as it exists in the persistence layer
// this is very similar to models.Resource, but COULD vary, so we should
// keep them separated
type Resource struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	SourceID string `json:"source_id"`
}
