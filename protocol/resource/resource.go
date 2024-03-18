package resource

type Metadata struct {
	From      string `json:"from"`
	Kind      string `json:"kind"`
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}
