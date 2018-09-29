package models

// Order - Strict order schema
type Order struct {
	// Order id
	ID string `json:"id"`

	// Courier id
	CourierID string `json:"courier_id"`

	// Order creation time in UTC format(ms)
	CreatedAt int64 `json:"created_at,omitempty"`

	// Order cancellation time in UTC format(ms)
	DeliveredAt int64 `json:"delivered_at"`

	Destination *Location `json:"destination,omitempty"`

	Source *Location `json:"source,omitempty"`
}
