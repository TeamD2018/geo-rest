// +build tarantool

package services

import (
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
	"time"
)

var (
	testName  = "Vasya"
	testPhone = "+79031109865"
	testLat   = 55.797344
	testLon   = 37.537746
)

type TarantoolRouteTestSuite struct {
	suite.Suite
	client          *tarantool.Connection
	routeDAO        *TarantoolRouteDAO
	ordersDao       *OrdersElasticDAO
	couriersDao     *CouriersElasticDAO
	pool            *dockertest.Pool
	resource        *dockertest.Resource
	logger          *zap.Logger
	testCourier     *models.Courier
	testOrderCreate *models.OrderCreate
}

func (s *TarantoolRouteTestSuite) BeforeTest(suiteName, testName string) {
	s.testCourier = &models.Courier{
		Name:  testName,
		ID:    uuid.NewV4().String(),
		Phone: &testPhone,
		Location: &models.Location{
			Point: elastic.GeoPointFromLatLon(testLat, testLon),
		},
	}
}

func (s *TarantoolRouteTestSuite) SetupSuite() {
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
	s.routeDAO = NewTarantoolRouteDAO(c, s.logger)
	//s.ordersDao = NewOrdersElasticDAO(s.client, s.logger, s.couriersDao, "")
	//s.couriersDao = NewCouriersElasticDAO(s.client, s.logger, "", DefaultCouriersReturnSize)
	err = migrations.Driver{Client: c, Logger: zap.NewExample()}.Run()

	if err != nil {
		log.Fatal(err)
		s.logger.Fatal("fail to perform migrations", zap.Error(err))
	}
}

func (s *TarantoolRouteTestSuite) TestCreateCourierOK() {
	err := s.routeDAO.CreateCourier(s.testCourier.ID)
	if !s.NoError(err) {
		return
	}
	err = s.routeDAO.DeleteCourier(s.testCourier.ID)
	if !s.NoError(err) {
		return
	}
}

func (s *TarantoolRouteTestSuite) TestAddPointToRouteOK() {
	err := s.routeDAO.CreateCourier(s.testCourier.ID)
	if !s.NoError(err) {
		return
	}
	err = s.routeDAO.AddPointToRoute(s.testCourier.ID, &models.PointWithTs{
		Point: s.testCourier.Location.Point, Ts: uint64(time.Now().Unix()),
	})
	if !s.NoError(err) {
		return
	}
}

func TestIntegrationTarantoolTestSuite(t *testing.T) {
	suite.Run(t, new(TarantoolRouteTestSuite))
}
