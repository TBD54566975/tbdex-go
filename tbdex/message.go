package tbdex

// MessageMetadata represents the metadata of a message e.g. RFQ, quote etc.
type MessageMetadata struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Kind       string `json:"kind"`
	ID         string `json:"id"`
	ExchangeID string `json:"exchangeId"`
	CreatedAt  string `json:"createdAt"`
	ExternalID string `json:"externalId,omitempty"`
	Protocol   string `json:"protocol"`
}
