package models

type OSMEntity struct {
	OSMID   int
	OSMType string
}

type Coordinates [][2]float64

type Geojson struct {
	Coordinates []Coordinates
}

type NominatimReverseResponse struct {
	Geojson *Geojson `json:"geojson"`
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
