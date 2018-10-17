package models

import "github.com/olivere/elastic"

type PointWithTs struct {
	Point *elastic.GeoPoint
	Ts    uint64
}
