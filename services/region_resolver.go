package services

import (
	"fmt"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/json-iterator/go"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

type NominatimRegionResolver struct {
	client     http.Client
	logger     *zap.Logger
	urlLookup  string
	urlReverse string
}

func (r *NominatimRegionResolver) Lookup(entity *models.OSMEntity) (string, error) {
	url := r.buildURLLookup(entity)
	resp, err := r.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	lookupResponse := models.NominatimLookupResponse{}
	if err := jsoniter.Unmarshal(body, &lookupResponse); err != nil {
		return "", err
	}
	return r.prettifyLookupResult(lookupResponse[0]), nil
}

func (r *NominatimRegionResolver) prettifyLookupResult(response *models.LookupResp) string {
	builder := strings.Builder{}
	if response.Address.County != "" {
		builder.WriteString(response.Address.County)
	}
	if response.Address.City != "" {
		builder.WriteString(response.Address.City)
	}
	builder.WriteString(", ")
	if response.Address.StateDistrict != "" {
		builder.WriteString(response.Address.StateDistrict)
	}
	builder.WriteString(", ")
	if response.Address.State != "" {
		builder.WriteString(response.Address.State)
	}
	return builder.String()
}

func (r *NominatimRegionResolver) ResolveRegion(entity *models.OSMEntity) (models.Polygon, error) {
	url := r.buildURLReverse(entity)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	regionResp := &models.NominatimReverseResponse{}
	if err := jsoniter.Unmarshal(body, regionResp); err != nil {
		return nil, err
	}
	polygon := make(models.Polygon, len(regionResp.Geojson.Coordinates[0]))
	for i, p := range regionResp.Geojson.Coordinates[0] {
		polygon[i] = elastic.GeoPointFromLatLon(p[1], p[0])
	}
	return polygon, nil
}

func NewNominatimRegionResolver(nominatimURL string, logger *zap.Logger) *NominatimRegionResolver {
	return &NominatimRegionResolver{
		client:     http.Client{},
		urlLookup:  nominatimURL + "/lookup?format=json&osm_ids=%s&accept-language=ru&osm_type=%s",
		urlReverse: nominatimURL + "/reverse?format=json&osm_id=%d&osm_type=%s&polygon_geojson=1&accept_language=ru",
		logger:     logger,
	}
}

func (r *NominatimRegionResolver) buildURLReverse(entity *models.OSMEntity) string {
	return fmt.Sprintf(r.urlReverse, entity.OSMID, entity.OSMType)
}

func (r *NominatimRegionResolver) buildURLLookup(entity *models.OSMEntity) string {
	return fmt.Sprintf(r.urlLookup, fmt.Sprintf("%s%d", entity.OSMType, entity.OSMID))
}
