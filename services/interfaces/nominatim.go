package interfaces

import "github.com/TeamD2018/geo-rest/models"

type IRegionResolver interface {
	ResolveRegion(entity *models.OSMEntity) (models.FlatPolygon, error)
}

type LookupInterface interface {
	Lookup(entity *models.OSMEntity) (string, error)
}
