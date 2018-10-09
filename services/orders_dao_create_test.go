// +build elastic

package services

import (
	"context"
	"fmt"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"testing"
)

type OrdersCreateTestSuite struct {
	suite.Suite
	client          *elastic.Client
	ordersDao       *OrdersElasticDAO
	couriersDao     *CouriersElasticDAO
	pool            *dockertest.Pool
	resource        *dockertest.Resource
	logger          *zap.Logger
	testCourier     *models.Courier
	testOrderCreate *models.OrderCreate
}

func TestIntegrationOrdersCreateSuite(t *testing.T) {
	suite.Run(t, new(OrdersCreateTestSuite))
}

func (s *OrdersCreateTestSuite) BeforeTest(suiteName, testName string) {
	s.ordersDao.EnsureMapping()
	s.couriersDao.EnsureMapping()
	testCourierCreate := models.CourierCreate{Name: "Test"}

	s.testCourier, _ = s.couriersDao.Create(&testCourierCreate)
	s.testOrderCreate = &models.OrderCreate{
		CourierID: &s.testCourier.ID,
		Destination: models.Location{
			Point: elastic.GeoPointFromLatLon(1, 1),
		},
		Source: models.Location{
			Point: elastic.GeoPointFromLatLon(1, 1),
		},
	}
	s.client.Refresh(s.ordersDao.Index, s.couriersDao.index).Do(context.Background())

}

func (s *OrdersCreateTestSuite) AfterTest(suiteName, testName string) {
	s.client.DeleteIndex(s.couriersDao.index, s.ordersDao.Index).Do(context.Background())
}

func (s *OrdersCreateTestSuite) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func (s *OrdersCreateTestSuite) SetupSuite() {
	log.SetFlags(log.Lshortfile)
	pool, err := dockertest.NewPool("")
	if err != nil {
		s.FailNow("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("bitnami/elasticsearch", "latest", []string{})
	if err != nil {
		s.FailNow("Could not start resource: %s", err)
	}

	var c *elastic.Client

	if err := pool.Retry(func() error {
		addr := fmt.Sprintf("http://localhost:%s", resource.GetPort("9200/tcp"))

		var err error
		c, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
		if err != nil {
			return err
		}

		_, _, err = c.Ping(addr).Do(context.Background())

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	s.client = c
	s.pool = pool
	s.resource = resource
	s.logger = zap.NewNop()
	s.couriersDao = NewCouriersElasticDAO(s.client, s.logger, "", DefaultCouriersReturnSize)
	s.ordersDao = NewOrdersElasticDAO(s.client, s.logger, s.couriersDao, "")
}

func (s OrdersCreateTestSuite) TestOrdersElasticDAO_EnsureMapping() {
	err := s.ordersDao.EnsureMapping()
	if !s.Assert().NoError(err) {
		return
	}

	exists, err := s.client.IndexExists(s.ordersDao.Index).Do(context.Background())
	s.Assert().NoError(err)
	s.Assert().True(exists)
}

func (s OrdersCreateTestSuite) TestOrdersElasticDAO_Create() {
	created, err := s.ordersDao.Create(s.testOrderCreate)
	if !s.Assert().NoError(err) {
		return
	}
	res, err := s.client.
		Exists().
		Index(s.ordersDao.Index).
		Type("_doc").
		Id(created.ID).
		Do(context.Background())
	s.Assert().NoError(err)
	s.Assert().True(res)
}

func (s OrdersCreateTestSuite) TestOrdersElasticDAO_NoCreateIfNoCourier() {
	cid := "550e8400-e29b-41d4-a716-446655440000"
	s.testOrderCreate.CourierID = &cid
	created, err := s.ordersDao.Create(s.testOrderCreate)
	s.Assert().Error(err)
	s.Assert().Nil(created)
}
