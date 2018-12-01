package services

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/TeamD2018/geo-rest/services/photon"
	"github.com/json-iterator/go"
	"strings"
)

type PhotonSuggestEngine struct {
	Tags    []string
	OSMType string
	Limit   int
}

type SuggestionFeatures struct {
	Features []*SuggestionFeature `json:"features"`
}

type SuggestionFeature struct {
	Properties Properties `json:"properties"`
}

type Properties struct {
	OSMID   int64  `json:"osm_id"`
	OSMType string `json:"osm_type"`
	Name    string `json:"name"`
	State   string `json:"state"`
	City    string `json:"city"`
	Country string `json:"country"`
}

func (p *Properties) String() string {
	const delimiter = ", "
	builder := strings.Builder{}
	builder.WriteString(p.Name)
	if p.City != "" {
		builder.WriteString(delimiter)
		builder.WriteString(p.City)
	}
	builder.WriteString(delimiter)
	builder.WriteString(p.State)
	builder.WriteString(delimiter)
	builder.WriteString(p.Country)
	return builder.String()
}

func (ops *PhotonSuggestEngine) ParseSearchResponse(response interface{}) interface{} {
	result := response.([]byte)
	features := SuggestionFeatures{}
	jsoniter.Unmarshal(result, &features)
	filteredFeatures := make([]*models.OSMPolygonSuggestion, 0, len(features.Features))
	for _, feature := range features.Features {
		props := feature.Properties
		if props.OSMType == ops.OSMType {
			filteredFeatures = append(filteredFeatures, &models.OSMPolygonSuggestion{
				OSMID:   props.OSMID,
				OSMType: props.OSMType,
				Name:    props.Name,
			})
		}
	}
	return filteredFeatures
}

func (ops *PhotonSuggestEngine) CreateSearchRequest(input string, ) (interface{}) {
	return photon.NewSearchQuery(input, ops.Limit, ops.Tags...)
}
