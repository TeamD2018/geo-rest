package services

import (
	"fmt"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/json-iterator/go"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type RegionResolver struct {
	client     http.Client
	logger     *zap.Logger
	urlLookup  string
	urlReverse string
}

func (r *RegionResolver) Lookup(entity *models.OSMEntity) (string, error) {
	url := r.buildURLLookup(entity)
	resp, err := r.client.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	lookupResponse := models.NominatimLookupResponse{}
	if err := jsoniter.Unmarshal(body, &lookupResponse); err != nil {
		return "", err
	}
	return lookupResponse[0].DisplayName, nil
}

func (r *RegionResolver) ResolveRegion(entity *models.OSMEntity) (*models.Polygon, error) {
	url := r.buildURLReverse(entity)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	regionResp := &models.NominatimReverseResponse{}
	if err := jsoniter.Unmarshal(body, regionResp); err != nil {
		return nil, err
	}
	polygon := models.NewPolygon(len(regionResp.Geojson.Coordinates[0]))
	for i, p := range regionResp.Geojson.Coordinates[0] {
		polygon.Points[i] = elastic.GeoPointFromLatLon(p[0], p[1])
	}
	return polygon, nil
}

func NewRegionResolver(nominatimURL string, logger *zap.Logger) *RegionResolver {
	return &RegionResolver{
		client:     http.Client{},
		urlLookup:  nominatimURL + "/lookup?format=json&osm_ids=%s&accept-language=ru&osm_type=%s",
		urlReverse: nominatimURL + "/reverse?format=json&osm_id=%d&osm_type=%s&polygon_geojson=1&accept_language=ru",
		logger:     logger,
	}
}

func (r *RegionResolver) buildURLReverse(entity *models.OSMEntity) string {
	return fmt.Sprintf(r.urlReverse, entity.OSMID, entity.OSMType)
}

func (r *RegionResolver) buildURLLookup(entity *models.OSMEntity) string {
	return fmt.Sprintf(r.urlLookup, fmt.Sprintf("%s%d", entity.OSMType, entity.OSMID))
}