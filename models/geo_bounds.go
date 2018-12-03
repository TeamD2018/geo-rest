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

type FlatPolygon []*elastic.GeoPoint
