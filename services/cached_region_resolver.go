package services

import (
	"errors"
	"github.com/TeamD2018/geo-rest/models"
)

type CachedRegionResolver struct {
	NominatimResolver *NominatimRegionResolver
	TarantoolResolver *TarantoolRegionResolver
}

func (r *CachedRegionResolver) ResolveRegion(entity *models.OSMEntity) (*models.Polygon, error) {
	if entity == nil {
		return nil, errors.New("entity is nil")
	}
	var polygon models.Polygon
	if polygon, err := r.TarantoolResolver.ResolveRegion(entity); err == models.ErrEntityNotFound {
		if polygon, err = r.NominatimResolver.ResolveRegion(entity); err != nil {
			return nil, err
		}
		if err := r.TarantoolResolver.SaveToCache(entity.OSMID, polygon); err != nil {
			return polygon, err
		}
	}
	return &polygon, nil
}
