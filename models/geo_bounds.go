package models

import "github.com/olivere/elastic"

type SquareField struct {
	HighLeftPoint  *elastic.GeoPoint
	DownRightPoint *elastic.GeoPoint
}

type CircleField struct {
	Center *elastic.GeoPoint
	Radius int
}
