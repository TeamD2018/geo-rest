package models

// Courier - Strict courier schema
type Courier struct {
	ID string `json:"id"`

	Name string `json:"name,omitempty"`

	Location *Location `json:"location,omitempty"`

	// Phone number in international format (without '+')
	Phone *string `json:"phone,omitempty"`

	// Time in UTC format(ms)
	LastSeen    *int64 `json:"last_seen,omitempty"`
	OrdersCount int    `json:"orders_count"`
	IsActive    bool   `json:"is_active,omitempty"`
}
