// +build tarantool

package services

import (
	"context"
	"fmt"
	"github.com/TeamD2018/geo-rest/migrations"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"log"
	"testing"
)

var (
	testAddress        = "Some Address"
	clearCacheFuncName = "clear_cache"
)

type TarantoolResolverTestSuite struct {
	suite.Suite
	client          *tarantool.Connection
	resolver        *TntResolver
	ordersDao       *OrdersElasticDAO
	couriersDao     *CouriersElasticDAO
	pool            *dockertest.Pool
	resource        *dockertest.Resource
	logger          *zap.Logger
	testCourier     *models.Courier
	testOrderCreate *models.OrderCreate
}

func (s *TarantoolResolverTestSuite) AfterTest(suiteName, testName string) {
	_, err := s.client.Call17(clearCacheFuncName, make([]interface{}, 0))
	if !s.NoError(err) {
		return
	}
}

func (s *TarantoolResolverTestSuite) BeforeTest(suiteName, testName string) {
	s.testCourier = &models.Courier{
		Name:  testName,
		ID:    uuid.NewV4().String(),
		Phone: &testPhone,
		Location: &models.Location{
			Point: elastic.GeoPointFromLatLon(testLat, testLon),
		},
	}
}

func (s *TarantoolResolverTestSuite) SetupSuite() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		s.FailNow("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("tarantool/tarantool", "2", []string{})
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
	s.resolver = NewTntResolver(c, s.logger)
	//s.ordersDao = NewOrdersElasticDAO(s.client, s.logger, s.couriersDao, "")
	//s.couriersDao = NewCouriersElasticDAO(s.client, s.logger, "", DefaultCouriersReturnSize)
	err = migrations.Driver{Client: c, Logger: zap.NewExample()}.Run()

	if err != nil {
		log.Fatal(err)
		s.logger.Fatal("fail to perform migrations", zap.Error(err))
	}
}

func (s *TarantoolResolverTestSuite) TestSaveToCache_OK() {
	location := models.Location{
		Point:   elastic.GeoPointFromLatLon(testLat, testLon),
		Address: &testAddress,
	}
	err := s.resolver.SaveToCache(&location)
	if !s.NoError(err) {
		return
	}
	var point = make([]interface{}, 0)
	err = s.client.Call17Typed(reversResolveFuncName, []interface{}{*location.Address}, &point)
	if !s.NoError(err) {
		return
	}
	if s.Len(point, 1) {
		return
	}
	if s.Contains(point, testAddress) {
		return
	}
}

func (s *TarantoolResolverTestSuite) TestRevertResolve_OK() {
	point := elastic.GeoPointFromLatLon(testLat, testLon)
	location := models.Location{
		Address: &testAddress,
		Point:   point,
	}
	err := s.resolver.SaveToCache(&location)
	if !s.NoError(err) {
		return
	}
	location.Point = nil

	err = s.resolver.Resolve(&location, context.Background())

	if !s.NoError(err) {
		return
	}

	if !s.NotNil(location.Point) {
		return
	}

	if !s.Equal(point, location.Point) {
		return
	}
}

func (s *TarantoolResolverTestSuite) TestResolve_OK() {
	point := elastic.GeoPointFromLatLon(testLat, testLon)
	location := models.Location{
		Address: &testAddress,
		Point:   point,
	}
	err := s.resolver.SaveToCache(&location)
	if !s.NoError(err) {
		return
	}
	location.Address = nil

	err = s.resolver.Resolve(&location, context.Background())

	if !s.NoError(err) {
		return
	}

	if !s.NotNil(location.Address) {
		return
	}

	if !s.Equal(testAddress, *location.Address) {
		return
	}
}

func (s *TarantoolResolverTestSuite) TestResolve_NotFound() {
	point := elastic.GeoPointFromLatLon(testLat, testLon)
	location := models.Location{
		Point:   point,
	}
	err := s.resolver.Resolve(&location, context.Background())

	if !s.IsType(models.ErrEntityNotFound, err) {
		return
	}
}

func (s *TarantoolResolverTestSuite) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func TestIntegrationTarantoolResolverTestSuite(t *testing.T) {
	suite.Run(t, new(TarantoolResolverTestSuite))
}
