package models

import "github.com/olivere/elastic"

type BoxField struct {
	TopLeftPoint     *elastic.GeoPoint
	BottomRightPoint *elastic.GeoPoint
}

type CircleField struct {
	Center *elastic.GeoPoint
	Radius int
}

type Polygon struct {
	Points []*elastic.GeoPoint
}

func NewPolygon(numPoints int) *Polygon {
	return &Polygon{make([]*elastic.GeoPoint, numPoints)}
}
