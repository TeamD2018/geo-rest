package openapi

// OrderBase - Basic courier fields
type OrderBase struct {

	// Order id
	Id int64 `json:"id,omitempty"`

	// Courier id
	CourierId int64 `json:"courier_id,omitempty"`

	// Order creation time in UTC format(ms)
	Created int64 `json:"created,omitempty"`

	// Order cancellation time in UTC format(ms)
	Done int64 `json:"done,omitempty"`

	Destination *Location `json:"destination,omitempty"`

	Soruce *Location `json:"soruce,omitempty"`
}
