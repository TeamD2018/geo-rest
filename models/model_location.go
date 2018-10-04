package models

import "github.com/olivere/elastic"

// Location - Geo-point (longitude, latitude)
type GeoPoint struct {
	// Latitude
	Lat float32 `json:"lat"`

	// Longitude
	Lon float32 `json:"lon"`
}

type Location struct {
	//Geopoint for location
	Point *elastic.GeoPoint `json:"point,omitempty"`

	//Address for location
	Address *string `json:"address,omitempty"`
}
