package parameters

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
)

// Location - Geo-point (longitude, latitude)
type CircleFieldQuery struct {
	Size int `form:"size"`
	// Latitude
	Lat float64 `form:"lat" binding:"min=-90,max=90"`

	// Longitude
	Lon float64 `form:"lon" binding:"min=-180,max=180"`

	ActiveOnly bool `form:"active_only"`
	Radius     int  `form:"radius" binding:"required,gt=0"`
}

type BoxFieldQuery struct {
	Size       int     `form:"size"`
	ActiveOnly bool    `form:"active_only"`
	TopLeftLat float64 `form:"top_left_lat" binding:"min=-90,max=90"`
	TopLeftLon float64 `form:"top_left_lon" binding:"min=-180,max=180"`

	BottomRightLat float64 `form:"bottom_right_lat" binding:"min=-90,max=90"`
	BottomRightLon float64 `form:"bottom_right_lon" binding:"min=-180,max=180"`
}

func (b *BoxFieldQuery) ToBoxField() *models.BoxField {
	boxField := &models.BoxField{
		TopLeftPoint:     elastic.GeoPointFromLatLon(b.TopLeftLat, b.TopLeftLon),
		BottomRightPoint: elastic.GeoPointFromLatLon(b.BottomRightLat, b.BottomRightLon),
	}
	return boxField
}

func (c *CircleFieldQuery) ToCircleField() *models.CircleField {
	circleField := &models.CircleField{
		Center: elastic.GeoPointFromLatLon(c.Lat, c.Lon),
		Radius: c.Radius,
	}
	return circleField
}
