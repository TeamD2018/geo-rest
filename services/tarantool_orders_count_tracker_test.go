// +build tarantool

package services

import (
	"fmt"
	"github.com/TeamD2018/geo-rest/migrations"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"log"
	"testing"
)

var (
	courierTestID = "550e8400-e29b-41d4-a716-446655440000"
)

type TarantoolOrdersCountTrackerTestSuite struct {
	suite.Suite
	client   *tarantool.Connection
	pool     *dockertest.Pool
	resource *dockertest.Resource
	logger   *zap.Logger
	tracker  *TarantoolOrdersCountTracker
}

func (s *TarantoolOrdersCountTrackerTestSuite) SetupSuite() {
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
	s.logger = zap.NewExample()
	s.tracker = NewTarantoolOrdersCountTracker(s.client, s.logger)
	err = migrations.Driver{Client: c, Logger: zap.NewExample()}.Run()
	if err != nil {
		log.Fatal(err)
		s.logger.Fatal("fail to perform migrations", zap.Error(err))
	}
}

func (s *TarantoolOrdersCountTrackerTestSuite) AfterTest(suiteName, testName string) {
	s.NoError(s.tracker.Drop(courierTestID))
}

func (s *TarantoolOrdersCountTrackerTestSuite) TestTarantoolOrdersCountTracker_Inc_OK() {
	have := models.Couriers{{ID: courierTestID, OrdersCount: -1}}
	if !s.NoError(s.tracker.Inc(courierTestID)) {
		return
	}
	if !s.NoError(s.tracker.Sync(have)) {
		return
	}
	s.Equal(1, have[0].OrdersCount)
}

func (s *TarantoolOrdersCountTrackerTestSuite) TestTarantoolOrdersCountTracker_IncAndGet_OK() {
	current, err := s.tracker.IncAndGet(courierTestID)
	if !s.NoError(err) {
		return
	}
	s.Equal(1, current)
}

func (s *TarantoolOrdersCountTrackerTestSuite) TestTarantoolOrdersCountTracker_DecAndGet_OK() {
	current, err := s.tracker.DecAndGet(courierTestID)
	if !s.NoError(err) {
		return
	}
	s.Equal(0, current)
}

func (s *TarantoolOrdersCountTrackerTestSuite) TestTarantoolOrdersCountTracker_Inc2xDecAndGet_OK() {
	if !s.NoError(s.tracker.Inc(courierTestID)) {
		return
	}
	if !s.NoError(s.tracker.Inc(courierTestID)) {
		return
	}
	current, err := s.tracker.DecAndGet(courierTestID)
	if !s.NoError(err) {
		return
	}
	s.Equal(1, current)
}

func (s *TarantoolOrdersCountTrackerTestSuite) TestTarantoolOrdersCountTracker_Dec_OK() {
	have := models.Couriers{{ID: courierTestID, OrdersCount: -1}}

	if !s.NoError(s.tracker.Dec(courierTestID)) {
		return
	}
	if !s.NoError(s.tracker.Sync(have)) {
		return
	}
	s.Zero(have[0].OrdersCount)
}

func (s *TarantoolOrdersCountTrackerTestSuite) TestTarantoolOrdersCountTracker_Inc2x_OK() {
	have := models.Couriers{{ID: courierTestID, OrdersCount: -1}}

	if !s.NoError(s.tracker.Inc(courierTestID)) {
		return
	}
	if !s.NoError(s.tracker.Inc(courierTestID)) {
		return
	}
	if !s.NoError(s.tracker.Sync(have)) {
		return
	}
	s.Equal(2, have[0].OrdersCount)
}
func (s *TarantoolOrdersCountTrackerTestSuite) TestTarantoolOrdersCountTracker_Inc2xDec1x_OK() {
	have := models.Couriers{{ID: courierTestID, OrdersCount: -1}}

	if !s.NoError(s.tracker.Inc(courierTestID)) {
		return
	}
	if !s.NoError(s.tracker.Inc(courierTestID)) {
		return
	}

	if !s.NoError(s.tracker.Dec(courierTestID)) {
		return
	}

	if !s.NoError(s.tracker.Sync(have)) {
		return
	}
	s.Equal(1, have[0].OrdersCount)
}

func (s *TarantoolOrdersCountTrackerTestSuite) TestTarantoolOrdersCountTracker_Sync_OK() {
	have := models.Couriers{{ID: courierTestID, OrdersCount: -1}}
	err := s.tracker.Sync(have)
	if !s.NoError(err) {
		return
	}
	s.Zero(have[0].OrdersCount)
}

func (s *TarantoolOrdersCountTrackerTestSuite) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func TestIntegrationTarantoolOrdersCountTrackerTestSuite(t *testing.T) {
	suite.Run(t, new(TarantoolOrdersCountTrackerTestSuite))
}
