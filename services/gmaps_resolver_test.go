//+build GOOGLE_TAKE_MY_MONEY

package services

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"log"
	"os"
	"testing"
)

func TestIntegrationGMapsResolver(t *testing.T) {
	suite.Run(t, new(GMapsResolverTestSuite))
}

type GMapsResolverTestSuite struct {
	suite.Suite
	client  *maps.Client
	service *GMapsResolver
	Logger  *zap.Logger
}

func (s *GMapsResolverTestSuite) SetupSuite() {
	log.SetFlags(log.Lshortfile)
	key := os.Getenv("GMAPS_API_KEY")
	client, err := maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		log.Fatal(err)
	}
	s.client = client
	s.Logger, _ = zap.NewDevelopment()
	s.service = NewGMapsResolver(client, s.Logger)
}

func (s GMapsResolverTestSuite) TestGMapsResolver_Resolve_Addr() {
	addr := "1600 Amphitheatre Pkwy, Mountain View, CA 94043, США"
	location := &models.Location{
		Point:   nil,
		Address: &addr,
	}
	err := s.service.Resolve(location)
	s.NoError(err)
	s.InDelta(37.4231778, location.Point.Lat, 0.001)
	s.InDelta(-122.0852514, location.Point.Lon, 0.001)
}

func (s GMapsResolverTestSuite) TestGMapsResolver_Resolve_Point() {
	addr := "1600 Amphitheatre Pkwy, Mountain View, CA 94043, США"
	location := &models.Location{
		Point: elastic.GeoPointFromLatLon(37.4231778, -122.0852514),
	}
	err := s.service.Resolve(location)
	s.NoError(err)
	s.EqualValues(addr, *location.Address)
}
