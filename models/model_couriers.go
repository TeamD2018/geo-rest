package models

type Couriers []*Courier

type CouriersResponse struct {
	Couriers Couriers    `json:"couriers,omitempty"`
	Polygon  FlatPolygon `json:"polygon,omitempty"`
}
