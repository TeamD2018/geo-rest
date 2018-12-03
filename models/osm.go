package models

import (
	"fmt"
	"encoding/json"
)

type OSMEntity struct {
	OSMID   int    `json:"osm_id"`
	OSMType string `json:"osm_type"`
}

type Address struct {
	StateDistrict string `json:"state_district"`
	City          string `json:"city"`
	State         string `json:"state"`
	County        string `json:"county"`
}

type LookupResp struct {
	DisplayName string   `json:"display_name"`
	Address     *Address `json:"address"`
}

type NominatimLookupResponse []*LookupResp

type NominatimReverseResponse struct {
	Geojson *Geojson `json:"geojson"`
}

type Coordinates [][2]float64

type Geojson struct {
	Coordinates GeoJSONPolygon `json:"coordinates"`
}

type GeoJSONPolygon []Coordinates
type GeoJSONMultiPolygon []GeoJSONPolygon

type GeoJSON struct {
	Geojson *struct {
		Type string `json:"type"`
	} `json:"geojson"`
}

type GeoJSONWithPolygon struct {
	Geojson *struct {
		Coordinates GeoJSONPolygon `json:"coordinates"`
	} `json:"geojson"`
}

type GeoJSONWithMultiPolygon struct {
	Geojson *struct {
		Coordinates []GeoJSONPolygon `json:"coordinates"`
	} `json:"geojson"`
}

func (r *NominatimReverseResponse) UnmarshalJSON(b []byte) error {
	var t GeoJSON
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	r.Geojson = &Geojson{}
	if t.Geojson.Type == "Polygon" {
		var polygon GeoJSONWithPolygon
		if err := json.Unmarshal(b, &polygon); err != nil {
			fmt.Println("fail on polygon")
			return err
		}
		r.Geojson.Coordinates = polygon.Geojson.Coordinates
		return nil
	}
	if t.Geojson.Type == "MultiPolygon" {
		var multiPolygon GeoJSONWithMultiPolygon
		if err := json.Unmarshal(b, &multiPolygon); err != nil {
			return err
		}
		r.Geojson.Coordinates = multiPolygon.Geojson.Coordinates[0]
		return nil
	}
	return fmt.Errorf("unsupported polygon type")
}
