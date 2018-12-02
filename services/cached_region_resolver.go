package services

import (
	"errors"
	"github.com/TeamD2018/geo-rest/models"
)

type CachedRegionResolver struct {
	NominatimResolver *NominatimRegionResolver
	TarantoolResolver *TarantoolRegionResolver
}

func (r *CachedRegionResolver) ResolveRegion(entity *models.OSMEntity) (polygon models.Polygon, err error) {
	if entity == nil {
		return nil, errors.New("entity is nil")
	}
	if polygon, err = r.TarantoolResolver.ResolveRegion(entity); err == models.ErrEntityNotFound {
		if polygon, err = r.NominatimResolver.ResolveRegion(entity); err != nil {
			return nil, err
		}
		if err := r.TarantoolResolver.SaveToCache(entity.OSMID, polygon); err != nil {
			return polygon, err
		}
	}
	return polygon, nil
}