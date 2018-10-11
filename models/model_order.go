package models

// Order - Strict order schema
type Order struct {
	// Order id
	ID string `json:"id,omitempty"`

	// Courier id
	CourierID string `json:"courier_id"`

	// Order creation time in UTC format(ms)
	CreatedAt int64 `json:"created_at,omitempty"`

	// Order cancellation time in UTC format(ms)
	DeliveredAt int64 `json:"delivered_at,omitempty"`

	Destination Location `json:"destination,omitempty"`

	Source Location `json:"source,omitempty"`
}

type OrderCreate struct {
	CourierID   *string  `json:"courier_id"`
	Destination Location `json:"destination"`
	Source      Location `json:"source"`
}

type OrderUpdate struct {
	// Order id
	ID *string `json:"id,omitempty"`

	// Courier id
	CourierID *string `json:"courier_id,omitempty"`

	// Order creation time in UTC format(ms)
	CreatedAt *int64 `json:"created_at,omitempty"`

	// Order cancellation time in UTC format(ms)
	DeliveredAt *int64 `json:"delivered_at,omitempty"`

	//Destination of order
	Destination *Location `json:"destination,omitempty"`

	//Source of order
	Source *Location `json:"source,omitempty"`
}
