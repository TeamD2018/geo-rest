package services

import (
	"errors"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

const (
	saveToCacheRegionFuncName = "save_to_cache_region"
	regionResolveFuncName = "resolve_region"
)

type TarantoolRegionResolver struct {
	c *tarantool.Connection
	l *zap.Logger
}

func NewTarantoolRegionResolver(client *tarantool.Connection, logger *zap.Logger) *TarantoolRegionResolver {
	return &TarantoolRegionResolver{
		c: client,
		l: logger,
	}
}

func (t *TarantoolRegionResolver) ResolveRegion(entity *models.OSMEntity) (models.Polygon, error) {
	var polygonInterface = make([]interface{}, 0)
	err := t.c.Call17Typed(regionResolveFuncName, tarantool.IntKey{I: entity.OSMID}, &polygonInterface)
	if err != nil {
		t.l.Debug("tarantool return", zap.Error(err))
		return nil, err
	}
	if len(polygonInterface) == 0 {
		return nil, models.ErrEntityNotFound
	}
	t.l.Sugar().Debugw("asd", "resp", polygonInterface)
	polygonInterface = polygonInterface[0].([]interface{})
	polygonFromTnt := polygonInterface[1].([]interface{})
	polygon := make(models.Polygon, len(polygonFromTnt))
	for i, p := range polygonFromTnt {
		pp := p.([]interface{})
		polygon[i] = elastic.GeoPointFromLatLon(pp[0].(float64), pp[1].(float64))
	}
	return polygon, nil
}

func (t *TarantoolRegionResolver) SaveToCache(osmID int, polygon models.Polygon) error {
	if polygon == nil {
		return errors.New("nil polygon")
	}
	polygonTnt := make([]interface{}, 0)
	for _, p := range polygon {
		polygonTnt = append(polygonTnt, [2]float64{p.Lat, p.Lon})
	}
	_, err := t.c.Call17(saveToCacheRegionFuncName, []interface{}{osmID, polygonTnt})
	if err != nil {
		return err
	}
	return nil
}