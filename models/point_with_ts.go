package models

import "github.com/olivere/elastic"

type PointWithTs struct {
	Point *elastic.GeoPoint `json:"point"`
	Ts    uint64            `json:"timestamp"`
}

type Points []*PointWithTs
