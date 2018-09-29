package models

// Order - Strict order schema
type Order struct {
	// Order id
	ID string `json:"id"`

	// Courier id
	CourierID string `json:"courier_id"`

	// Order creation time in UTC format(ms)
	Created int64 `json:"created,omitempty"`

	// Order cancellation time in UTC format(ms)
	EndTime int64 `json:"end_time"`

	Destination *Location `json:"destination,omitempty"`

	Source *Location `json:"source,omitempty"`
}
