package models

import "github.com/json-iterator/go"

type OSMEntity struct {
	OSMID   int
	OSMType string
}

type Coordinates [][2]float64

type Geojson struct {
	Coordinates [][]Coordinates
}

type NominatimReverseResponse struct {
	Geojson *Geojson `json:"geojson"`
}

func (r *NominatimReverseResponse) UnmarshalJSON(b []byte) error {
	raw := make(map[string]interface{}, 0)
	if err := jsoniter.Unmarshal(b, raw); err != nil {
		return nil
	}
	if raw["geojson"].(map[string]interface{})["type"] == "MultiPolygon" {
		r.Geojson = &Geojson{Coordinates: make([][]Coordinates, len(raw["geojson"].(map[string]interface{})["coordinates"].([]interface{})))}
		for i, pI := range raw["geojson"].(map[string]interface{})["coordinates"].([][]Coordinates) {
			r.Geojson.Coordinates[i] = pI
		}
	} else {
		r.Geojson = &Geojson{Coordinates: [][]Coordinates{raw["geojson"].(map[string]interface{})["coordinates"].([]Coordinates)}}
	}
	return nil
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
