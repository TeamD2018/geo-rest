package models

type CourierUpdate struct {
	ID       *string   `json:"id"`
	Name     *string   `json:"name"`
	Location *Location `json:"location"`
	Phone    *string   `json:"phone"`
	LastSeen *int64    `json:"last_seen"`
}

type CourierCreate struct {
	Name  string  `json:"name" binding:"required"`
	Phone *string `json:"phone"`
}
