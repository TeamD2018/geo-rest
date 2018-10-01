// +build elastic

package services

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"testing"
)


type OrdersCreateSuite struct {
	suite.Suite
	client   *elastic.Client
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (s *OrdersCreateSuite) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func (s *OrdersCreateSuite) SetupSuite() {
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
}

func (s OrdersCreateSuite) TestOrdersElasticDAO_EnsureMapping() {
	logger, _ := zap.NewDevelopment()

	c := NewOrdersElasticDAO(s.client, logger, OrdersIndex)

	err := c.EnsureMapping()
	s.Assert().NoError(err)

	exists, err := s.client.IndexExists(OrdersIndex).Do(context.Background())
	s.Assert().NoError(err)
	s.Assert().True(exists)
}

func TestIntegrationOrdersSuite(t *testing.T) {
	suite.Run(t, new(OrdersCreateSuite))
}