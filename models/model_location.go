package models

import "github.com/olivere/elastic"

type Location struct {
	//Geopoint for location
	Point *elastic.GeoPoint `json:"point,omitempty"`

	//Address for location
	Address *string `json:"address,omitempty"`
}
