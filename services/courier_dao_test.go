// +build elastic

package services

import (
	"context"
	"fmt"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type CourierTestSuite struct {
	suite.Suite
	client   *elastic.Client
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (s *CourierTestSuite) AfterTest(suiteName, testName string) {
	if contains(testsWithDeleteIndex, testName) {
		s.DeleteIndex()
	}
}

var (
	testsWithCreateIndex = []string{
		"TestCreateCourierWithNameAndPhone",
		"TestCreateCourierWithName",
	}
	testsWithDeleteIndex = []string {
		"TestCreateCourierWithNameAndPhone",
		"TestCreateCourierWithName",
		"TestCouriersElasticDAO_EnsureMapping",
	}
)

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (s *CourierTestSuite) BeforeTest(suiteName, testName string) {
	if contains(testsWithCreateIndex, testName) {
		s.CreateIndex()
	}
}

func (s *CourierTestSuite) CreateIndex() {
	index, mapping := s.GetService().GetMapping()
	_, err := s.client.CreateIndex(index).BodyString(mapping).Do(context.Background())
	s.Assert().NoError(err)
}

func (s *CourierTestSuite) DeleteIndex() {
	index, _ := s.GetService().GetMapping()
	_, err := s.client.DeleteIndex(index).Do(context.Background())
	s.Assert().NoError(err)
}

func (s *CourierTestSuite) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func (s *CourierTestSuite) SetupSuite() {
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
}

func (s *CourierTestSuite) GetService() *CouriersElasticDAO {
	return NewCouriersDAO(s.client, CourierIndex, nil)
}

func (s *CourierTestSuite) ClearCouriersFromElastic(couriersIDs ...string) error {
	for _, id := range couriersIDs {
		if _, err := s.client.Delete().Index(CourierIndex).Id(id).Do(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

func (s *CourierTestSuite) TestCreateCourierWithNameAndPhone() {
	service := s.GetService()

	phone := "79123456789"
	name := "Vasya"
	courier := &models.CourierCreate{
		Name:  name,
		Phone: &phone,
	}
	createdCourier, err := service.Create(courier)

	s.Assert().NoError(err)
	s.Assert().Equal(createdCourier.Name, name)
	s.Assert().Equal(*createdCourier.Phone, phone)

	s.ClearCouriersFromElastic(createdCourier.ID)
}

func (s *CourierTestSuite) TestCreateCourierWithName() {
	service := s.GetService()

	name := "Vasya"
	courier := &models.CourierCreate{
		Name: name,
	}
	createdCourier, err := service.Create(courier)

	s.Assert().NoError(err)
	s.Assert().Equal(createdCourier.Name, name)

	s.ClearCouriersFromElastic(createdCourier.ID)
}

func (s *CourierTestSuite) TestCouriersElasticDAO_EnsureMapping() {
	service := s.GetService()

	err := service.EnsureMapping()
	s.Assert().NoError(err)

	exists, err := s.client.IndexExists(CourierIndex).Do(context.Background())
	s.Assert().NoError(err)
	s.Assert().True(exists)
}

func TestIntegrationCouriersDAO(t *testing.T) {
	suite.Run(t, new(CourierTestSuite))
}
