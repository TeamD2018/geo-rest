package services

import (
	"context"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"log"
)

const (
	spaceGeoCacheName   = "geo_cache"
	indexName           = "address"
	saveToCacheFuncName = "save_to_cache"
	resolveFuncName     = "resolve"
)

type TntResolver struct {
	client *tarantool.Connection
}

func NewTntResolver(client *tarantool.Connection, logger *zap.Logger) *TntResolver {
	return &TntResolver{client: client}
}

func (tnt *TntResolver) Resolve(location *models.Location, ctx context.Context) error {
	var point = make([]interface{}, 0)
	if err := tnt.client.GetTyped(spaceGeoCacheName, indexName, tarantool.StringKey{*location.Address}, &point); err != nil {
		log.Println(err)
	} else {
		if len(point) == 0 {
			return models.ErrEntityNotFound
		}
		locFromTnt := point[1].(map[interface{}]interface{})
		location.Point = elastic.GeoPointFromLatLon(locFromTnt["lat"].(float64), locFromTnt["lon"].(float64))
	}
	return nil
}

func (tnt *TntResolver) SaveToCache(location *models.Location) error {
	if location == nil {
		return errors.New("location is nil")
	}
	point := map[string]float64{"lat": location.Point.Lat, "lon": location.Point.Lon}
	resp, err := tnt.client.Call17(saveToCacheFuncName, []interface{}{*location.Address, point})
	if err != nil {
		return err
	}
	log.Println(resp.Data)
	return nil
}
