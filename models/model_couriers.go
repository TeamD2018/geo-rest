package models

type Couriers []*Courier

type CouriersResponse struct {
	Couriers Couriers `json:"couriers"`
	Polygon  Polygon  `json:"polygon,omitempty"`
}
