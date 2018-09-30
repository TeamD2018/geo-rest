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

type SimpleTestSuite struct {
	suite.Suite
	client   *elastic.Client
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (s *SimpleTestSuite) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func (s *SimpleTestSuite) SetupSuite() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		s.FailNow("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("bitnami/elasticsearch", "6.4.1", []string{})
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

		_, _, err = client.Ping(addr).Do(context.Background())

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	s.client = c
	s.pool = pool
	s.resource = resource
}

func (s *SimpleTestSuite) TestCreateCourier() {
	service := NewCouriersDAO(s.client, CourierIndex, nil)
	phone := "+79123456789"
	name := "Vasya"
	courier := &models.CourierCreate{
		Name:  name,
		Phone: &phone,
	}
	createdCourier, err := service.Create(courier)
	s.Nil(err)
	s.Equal(createdCourier.Name, name)
}

func TestSimpleTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleTestSuite))
}
