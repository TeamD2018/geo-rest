package models

type Suggestion struct {
	Couriers Couriers `json:"couriers"`
	Orders   Orders   `json:"orders"`
}
