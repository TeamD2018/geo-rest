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
	GeoPoint *elastic.GeoPoint `json:"point"`

	//Address for location
	Address *string `json:"address,omitempty"`
}
