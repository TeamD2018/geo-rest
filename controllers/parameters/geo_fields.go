package parameters

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
)

// Location - Geo-point (longitude, latitude)
type CircleFieldQuery struct {
	// Latitude
	Lat float64 `form:"lat" validate:"min=-90,max=90"`

	// Longitude
	Lon float64 `form:"lon" validate:"min-180,max=180"`

	Radius int `form:"radius" validate:"gt=0"`
}

func (c *CircleFieldQuery) ToCircleField() *models.CircleField {
	circleField := &models.CircleField{
		Center: elastic.GeoPointFromLatLon(c.Lat, c.Lon),
		Radius: c.Radius,
	}
	return circleField
}
