package openapi

// CourierBase - Basic courier properties
type CourierBase struct {

	Id int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Location *Location `json:"location,omitempty"`

	// Phone number in international format (without '+')
	Phone string `json:"phone,omitempty"`

	// Time in UTC format(ms)
	LastSeen int64 `json:"last_seen,omitempty"`
}
