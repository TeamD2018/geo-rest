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
	spaceGeoCacheName     = "geo_cache"
	indexName             = "address"
	saveToCacheFuncName   = "save_to_cache"
	reversResolveFuncName = "revers_resolve"
	resolveFuncName       = "resolve"
)

type TntResolver struct {
	client *tarantool.Connection
}

func NewTntResolver(client *tarantool.Connection, logger *zap.Logger) *TntResolver {
	return &TntResolver{client: client}
}

func (tnt *TntResolver) reverseResolve(location *models.Location) error {
	var point = make([]interface{}, 0)
	if err := tnt.client.Call17Typed(reversResolveFuncName, tarantool.StringKey{S: *location.Address}, &point); err != nil {
		log.Println(err)
	} else {
		if len(point) == 0 {
			return models.ErrEntityNotFound
		}
		point = point[0].([]interface{})
		locFromTnt := point[1].([]interface{})
		location.Point = elastic.GeoPointFromLatLon(locFromTnt[0].(float64), locFromTnt[1].(float64))
	}
	return nil
}

func (tnt *TntResolver) resolve(location *models.Location) error {
	var address = make([]interface{}, 0)
	point := []float64{location.Point.Lat, location.Point.Lon}
	if err := tnt.client.Call17Typed(resolveFuncName, []interface{}{point}, &address); err != nil {
		log.Println(err)
		return err
	} else {
		log.Printf("%#v\n", address)
		if len(address) == 0 {
			return models.ErrEntityNotFound
		}
		addressStr := ((address[0].([]interface{})[0].([]interface{}))[0]).(string)
		location.Address = &addressStr
	}
	return nil
}

func (tnt *TntResolver) Resolve(location *models.Location, ctx context.Context) error {
	if location == nil {
		return errors.New("location is nil")
	}
	if location.Address == nil {
		return tnt.resolve(location)
	}
	return tnt.reverseResolve(location)
}

func (tnt *TntResolver) SaveToCache(location *models.Location) error {
	if location == nil {
		return errors.New("location is nil")
	}
	point := []float64{location.Point.Lat, location.Point.Lon}
	_, err := tnt.client.Call17(saveToCacheFuncName, []interface{}{*location.Address, point})
	if err != nil {
		return err
	}
	return nil
}
