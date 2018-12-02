package interfaces

import "github.com/TeamD2018/geo-rest/models"

type IRegionResolver interface {
	ResolveRegion(entity *models.OSMEntity) (models.Polygon, error)
}

type LookupInterface interface {
	Lookup(entity *models.OSMEntity) (string, error)
}