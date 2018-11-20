package services

import (
	"context"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
)

type GMapsResolver struct {
	Maps   *maps.Client
	Logger *zap.Logger
}

func NewGMapsResolver(client *maps.Client, logger *zap.Logger) *GMapsResolver {
	return &GMapsResolver{
		Maps:   client,
		Logger: logger,
	}
}

func (gm *GMapsResolver) Resolve(location *models.Location, ctx context.Context) error {
	if location == nil || (location.Point != nil && location.Address != nil) {
		return nil
	}
	if location.Address != nil {
		resolvedPoint, err := gm.resolveAddr(*location.Address, ctx)
		if err != nil {
			return err
		}
		location.Point = resolvedPoint
	}
	if location.Point != nil {
		resolvedAddr, err := gm.resolvePoint(location.Point, ctx)
		if err != nil {
			return err
		}
		location.Address = &resolvedAddr
		return nil
	}
	return nil
}
func (gm *GMapsResolver) resolvePoint(point *elastic.GeoPoint, ctx context.Context) (string, error) {
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: point.Lat,
			Lng: point.Lon,
		},
		ResultType: []string{
			"street_address",
			"administrative_area_level_2",
			"administrative_area_level_3",
			"colloquial_area",
			"sublocality",
			"premise",
			"subpremise",
			"natural_feature",
			"airport",
			"park",
			"point_of_interest",
		},
		Language: "ru",
	}
	results, err := gm.Maps.ReverseGeocode(ctx, req)
	if err != nil {
		return "", err
	}
	reversed := results[0]
	return reversed.FormattedAddress, nil
}

func (gm *GMapsResolver) resolveAddr(address string, ctx context.Context) (*elastic.GeoPoint, error) {
	req := &maps.GeocodingRequest{
		Address: address,
	}
	results, err := gm.Maps.Geocode(ctx, req)
	if err != nil {
		return nil, err
	}
	point := results[0].Geometry.Location
	return elastic.GeoPointFromLatLon(point.Lat, point.Lng), nil
}
