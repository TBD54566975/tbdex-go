package tbdex

// ResourceMetadata represents the metadata of a resource e.g. offering, balance etc.
type ResourceMetadata struct {
	From      string `json:"from"`
	Kind      string `json:"kind"`
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Protocol  string `json:"protocol"`
}
