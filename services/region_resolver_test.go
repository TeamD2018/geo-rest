// +build local

package services

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

const nominatimURL = "http://95.163.181.132:8080"

const (
	testOSMID = 1255680
	testOSMType = "R"
)

type NominatimTestSuite struct {
	suite.Suite
	resolver    *RegionResolver
}

func (s *NominatimTestSuite) BeforeTest(suiteName, testName string) {
	s.resolver = NewRegionResolver(nominatimURL, zap.NewNop())
}

func (s *NominatimTestSuite) TestResolveOK() {
	osmEntity := models.OSMEntity{
		OSMID: testOSMID,
		OSMType: testOSMType,
	}
	polygon, err := s.resolver.ResolveRegion(&osmEntity)
	if !s.NoError(err) {
		return
	}

	if !s.NotNil(polygon) {
		return
	}
}

func (s *NominatimTestSuite) TestLookupOK() {
	osmEntity := models.OSMEntity{
		OSMID: testOSMID,
		OSMType: testOSMType,
	}
	address, err := s.resolver.Lookup(&osmEntity)
	if !s.NoError(err) {
		return
	}
	if !s.Equal("район Северное Измайлово, Восточный административный округ, Москва, Центральный федеральный округ, Россия", address) {
		return
	}
}

func TestIntegrationNominatimSuite(t *testing.T) {
	suite.Run(t, new(NominatimTestSuite))
}