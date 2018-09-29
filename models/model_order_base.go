package models

// OrderBase - Basic courier fields
type OrderBase struct {
	// Order id
	ID *string `json:"id,omitempty"`

	// Courier id
	CourierID *string `json:"courier_id,omitempty"`

	// Order creation time in UTC format(ms)
	Created *int64 `json:"created,omitempty"`

	// Order cancellation time in UTC format(ms)
	EndTime *int64 `json:"end_time,omitempty"`

	//Destination of order
	Destination *Location `json:"destination,omitempty"`

	//Source of order
	Source *Location `json:"source,omitempty"`
}
