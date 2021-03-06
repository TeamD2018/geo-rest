package services

import (
	"bytes"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

var IzamilovoSuburbTagResponse = bytes.NewBufferString(`{"features":[{"geometry":{"coordinates":[37.7841549,55.7966159],"type":"Point"},"type":"Feature","properties":{"osm_id":1319245,"osm_type":"R","extent":[37.7429678,55.8028665,37.8055777,55.7596111],"country":"Russia","osm_key":"place","osm_value":"suburb","name":"Izmaylovo District","state":"Moscow"}},{"geometry":{"coordinates":[37.8847818,55.7910108],"type":"Point"},"type":"Feature","properties":{"osm_id":4430011881,"osm_type":"N","country":"Russia","osm_key":"place","city":"Balashikha","osm_value":"suburb","postcode":"143915","name":"Новое Измайлово","state":"Moscow Oblast"}},{"geometry":{"coordinates":[37.8244387,55.7792553],"type":"Point"},"type":"Feature","properties":{"osm_id":4862025988,"osm_type":"N","country":"Russia","osm_key":"place","city":"Ivanovskoye District","osm_value":"suburb","postcode":"105568","name":"Южное Измайлово","state":"Moscow"}},{"geometry":{"coordinates":[37.8218434,55.7965807],"type":"Point"},"type":"Feature","properties":{"osm_id":1625648082,"osm_type":"N","country":"Russia","osm_key":"place","city":"Vostochnoye Izmaylovo District","osm_value":"suburb","postcode":"105203","name":"Vostochnoye Izmaylovo","state":"Moscow"}},{"geometry":{"coordinates":[37.8033479,55.8067811],"type":"Point"},"type":"Feature","properties":{"osm_id":1625648129,"osm_type":"N","country":"Russia","osm_key":"place","city":"Severnoye Izmaylovo District","osm_value":"suburb","postcode":"105215","name":"Severnoye Izmaylovo","state":"Moscow"}}]}`)

type PhotonSuggestEngineTestSuite struct {
	suite.Suite
	engine *PhotonSuggestEngine
}

type lookupServiceMock struct {
	mock.Mock
}

func (lsm lookupServiceMock) Lookup(entity *models.OSMEntity) (string, error) {
	args := lsm.Called(entity)
	return args.String(0), args.Error(1)
}

func (s *PhotonSuggestEngineTestSuite) SetupSuite() {
	s.engine = &PhotonSuggestEngine{
		OSMType: "R",
		Limit:   10,
		Tags:    []string{"boundary:administrative", "place:suburb"},
		ConcurrentLookupService: &ConcurrentLookupService{
			Concurrency: 5,
		},
	}
}

func (s *PhotonSuggestEngineTestSuite) TestPhotonSuggestEngine_ParseSearchResponse_OK() {
	s.setupLookupMock(5, "Example", nil)
	response := IzamilovoSuburbTagResponse.Bytes()
	result, err := s.engine.ParseSearchResponse(response)
	if !s.NoError(err) {
		return
	}
	if !s.IsType([]*models.OSMPolygonSuggestion{}, result,) {
		return
	}
	expected := models.OSMPolygonSuggestion{
		OSMID:   1319245,
		OSMType: "R",
		Name:    "Example",
	}
	s.EqualValues(result, []*models.OSMPolygonSuggestion{&expected})
}
func (s *PhotonSuggestEngineTestSuite) setupLookupMock(concurrency int, retstr string, reterr error) {
	lookupMock := new(lookupServiceMock)
	lookupMock.On("Lookup", mock.Anything).Return(retstr, reterr)
	s.engine.ConcurrentLookupService = &ConcurrentLookupService{
		Concurrency:   5,
		LookupService: lookupMock,
	}
}

func TestPhotonSuggestEngineTestSuite(t *testing.T) {
	suite.Run(t, new(PhotonSuggestEngineTestSuite))
}
