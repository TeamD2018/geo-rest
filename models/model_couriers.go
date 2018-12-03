package models

type Couriers []*Courier

type CouriersResponse struct {
	Couriers Couriers `json:"couriers,omitempty"`
	Polygon  Polygon  `json:"polygon,omitempty"`
}
