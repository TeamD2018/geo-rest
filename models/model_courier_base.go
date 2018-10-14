package models

type CourierUpdate struct {
	ID          *string               `json:"-"`
	Name        *string               `json:"name,omitempty"`
	Location    *Location             `json:"location,omitempty"`
	Phone       *string               `json:"phone,omitempty"`
	LastSeen    *int64                `json:"last_seen,omitempty"`
}

type CourierCreate struct {
	Name  string  `json:"name" binding:"required"`
	Phone *string `json:"phone"`
}
