// +build tarantool

package services

import (
	"fmt"
	"github.com/TeamD2018/geo-rest/migrations"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"log"
	"testing"
)

var (
	clearRegionCacheFuncName = "clear_cache_region"
	resolveRegionFuncName    = "resolve_region"
	testOSMID                = 1255680
)

type TarantoolRegionResolverTestSuite struct {
	suite.Suite
	client   *tarantool.Connection
	resolver *TarantoolRegionResolver
	pool     *dockertest.Pool
	resource *dockertest.Resource
	logger   *zap.Logger
}

func (s *TarantoolRegionResolverTestSuite) AfterTest(suiteName, testName string) {
	_, err := s.client.Call17(clearRegionCacheFuncName, make([]interface{}, 0))
	if !s.NoError(err) {
		return
	}
}

func (s *TarantoolRegionResolverTestSuite) SetupSuite() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		s.FailNow("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("tarantool/tarantool", "1.10.2", []string{})
	if err != nil {
		s.FailNow("Could not start resource: %s", err)
	}

	var c *tarantool.Connection

	if err := pool.Retry(func() error {
		addr := fmt.Sprintf("localhost:%s", resource.GetPort("3301/tcp"))

		var err error
		c, err = tarantool.Connect(addr, tarantool.Opts{})
		if err != nil {
			return err
		}
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	s.client = c
	s.pool = pool
	s.resource = resource
	s.logger = zap.NewNop()
	s.resolver = NewTarantoolRegionResolver(c, s.logger)
	err = migrations.Driver{Client: c, Logger: zap.NewExample()}.Run()

	if err != nil {
		log.Fatal(err)
		s.logger.Fatal("fail to perform migrations", zap.Error(err))
	}
}

func (s *TarantoolRegionResolverTestSuite) TestSaveToCache_OK() {
	polygon := models.FlatPolygon{
		elastic.GeoPointFromLatLon(56.514792, 36.375407),
		elastic.GeoPointFromLatLon(56.673754, 39.858289),
		elastic.GeoPointFromLatLon(54.692269, 38.979383),
		elastic.GeoPointFromLatLon(55.146880, 36.122938),
		elastic.GeoPointFromLatLon(56.514792, 36.375407),
	}
	err := s.resolver.SaveToCache(testOSMID, polygon)
	if !s.NoError(err) {
		return
	}
	var result = make([]interface{}, 0)
	err = s.client.Call17Typed(resolveRegionFuncName, []interface{}{testOSMID}, &result)
	if !s.NoError(err) {
		return
	}
	if s.Len(result, 1) {
		return
	}
	if s.Contains(result, testAddress) {
		return
	}
}

func (s *TarantoolRegionResolverTestSuite) TestResolve_OK() {
	polygon := models.FlatPolygon{
		elastic.GeoPointFromLatLon(56.514792, 36.375407),
		elastic.GeoPointFromLatLon(56.673754, 39.858289),
		elastic.GeoPointFromLatLon(54.692269, 38.979383),
		elastic.GeoPointFromLatLon(55.146880, 36.122938),
		elastic.GeoPointFromLatLon(56.514792, 36.375407),
	}
	err := s.resolver.SaveToCache(testOSMID, polygon)
	if !s.NoError(err) {
		return
	}

	polygon, err = s.resolver.ResolveRegion(&models.OSMEntity{OSMID: testOSMID})

	if !s.NoError(err) {
		return
	}

	if !s.NotNil(polygon) {
		return
	}

	if !s.Len(polygon, 5) {
		return
	}
}

func (s *TarantoolRegionResolverTestSuite) TestResolve_NotFound() {
	_, err := s.resolver.ResolveRegion(&models.OSMEntity{OSMID: testOSMID})

	if !s.IsType(models.ErrEntityNotFound, err) {
		return
	}
}

func (s *TarantoolRegionResolverTestSuite) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func TestIntegrationTarantoolResolverRegionTestSuite(t *testing.T) {
	suite.Run(t, new(TarantoolRegionResolverTestSuite))
}
