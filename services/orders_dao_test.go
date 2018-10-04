// +build elastic

package services

import (
	"context"
	"fmt"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"testing"
)

type OrdersTestSuite struct {
	suite.Suite
	client      *elastic.Client
	ordersDao   *OrdersElasticDAO
	couriersDao *CouriersElasticDAO
	pool        *dockertest.Pool
	resource    *dockertest.Resource
	logger      *zap.Logger
	testCourier *models.Courier
	testOrder   *models.Order
}

func (s *OrdersTestSuite) BeforeTest(suiteName, testName string) {
	s.couriersDao.index = uuid.NewV4().String()
	s.couriersDao.EnsureMapping()
	s.ordersDao.Index = uuid.NewV4().String()
	s.ordersDao.EnsureMapping()
	testCourierCreate := models.CourierCreate{Name: "Test"}
	s.testCourier, _ = s.couriersDao.Create(&testCourierCreate)
	testOrderCreate := models.OrderCreate{
		CourierID: &s.testCourier.ID,
		Destination: models.Location{
			Point: elastic.GeoPointFromLatLon(1, 1),
		},
		Source: models.Location{
			Point: elastic.GeoPointFromLatLon(1, 1),
		},
	}
	s.testOrder, _ = s.ordersDao.Create(&testOrderCreate)
	s.client.Refresh(s.ordersDao.Index, s.couriersDao.index).Do(context.Background())
}

func (s *OrdersTestSuite) AfterTest(suiteName, testName string) {
	s.client.DeleteIndex(s.ordersDao.Index, s.couriersDao.index).Do(context.Background())
}

func (s *OrdersTestSuite) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func (s *OrdersTestSuite) SetupSuite() {
	log.SetFlags(log.Lshortfile)
	pool, err := dockertest.NewPool("")
	if err != nil {
		s.FailNow("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("bitnami/elasticsearch", "latest", []string{})
	if err != nil {
		log.Println(err)
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
	s.couriersDao = NewCouriersElasticDAO(s.client, nil, "")
	s.ordersDao = NewOrdersElasticDAO(s.client, nil, s.couriersDao, "")

}

func (s OrdersTestSuite) TestOrdersDao_Get() {
	got, err := s.ordersDao.Get(s.testOrder.ID)
	s.Assert().NoError(err)
	s.Assert().EqualValues(*s.testOrder, *got)
}

func (s OrdersTestSuite) TestOrdersElasticDAO_GetOrdersForCourier() {
	orders, err := s.ordersDao.GetOrdersForCourier(s.testCourier.ID, s.testOrder.CreatedAt, false)
	s.Assert().NoError(err)
	if s.Assert().Len(orders, 1) {
		s.Assert().Equal(orders[0].ID, s.testOrder.ID)
	}
}

func (s OrdersTestSuite) TestOrdersElasticDAO_GetOrdersForCourier_NoCourier() {
	orders, err := s.ordersDao.GetOrdersForCourier("550e8400-e29b-41d4-a716-446655440000", s.testOrder.CreatedAt, false)
	s.Assert().NoError(err)
	s.Assert().Empty(orders)
}

func (s OrdersTestSuite) TestOrdersElasticDAO_GetOrdersForCourier_TimeTreshold() {
	orders, err := s.ordersDao.GetOrdersForCourier(s.testCourier.ID, 0, false)
	s.Assert().NoError(err)
	s.Assert().Empty(orders)
}

func (s OrdersTestSuite) TestOrdersElasticDAO_GetOrdersForCourier_Desc() {
	orders, err := s.ordersDao.GetOrdersForCourier(s.testCourier.ID, s.testOrder.CreatedAt-10, true)
	s.Assert().NoError(err)
	if s.Assert().Len(orders, 1) {
		s.Assert().Equal(orders[0].ID, s.testOrder.ID)
	}
}

func (s OrdersTestSuite) TestOrdersElasticDAO_Update_OK() {
	update := models.OrderUpdate{
		ID:     &s.testOrder.ID,
		Source: &models.Location{Point: elastic.GeoPointFromLatLon(0, 0)},
	}
	expected := models.Order{
		ID:          s.testOrder.ID,
		CourierID:   s.testOrder.CourierID,
		Source:      *update.Source,
		Destination: s.testOrder.Destination,
		CreatedAt:   s.testOrder.CreatedAt,
	}
	updated, err := s.ordersDao.Update(&update)
	s.Assert().NoError(err)
	s.Assert().EqualValues(&expected, updated)
}

func (s OrdersTestSuite) TestOrdersElasticDAO_Delete_OK() {
	err := s.ordersDao.Delete(s.testOrder.ID)
	s.Assert().NoError(err)
	exists, err := s.client.Exists().Index(s.ordersDao.Index).Type("_doc").Id(s.testOrder.ID).Do(context.Background())
	s.Assert().NoError(err)
	s.Assert().False(exists)
}

func TestIntegrationOrdersSuite(t *testing.T) {
	suite.Run(t, new(OrdersTestSuite))
}
