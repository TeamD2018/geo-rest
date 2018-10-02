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
		"TestGetCourierByIDOK",
		"TestGetCourierByIDBadID",
		"TestUpdateCourierWithoutLocationOK",
		"TestUpdateCourierWithLocationOK",
		"TestGetCourierByID",
	}
	testsWithDeleteIndex = []string{
		"TestCreateCourierWithNameAndPhone",
		"TestCreateCourierWithName",
		"TestGetCourierByIDOK",
		"TestGetCourierByIDBadID",
		"TestUpdateCourierWithoutLocationOK",
		"TestUpdateCourierWithLocationOK",
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

func (s *CourierTestSuite) CreateCourier(courier *models.CourierCreate) string {
	service := s.GetService()
	resp, err := service.Create(courier)
	if !s.Assert().NoError(err) {
		s.Assert().Fail(err.Error())
	}
	return resp.ID
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
	return NewCouriersElasticDAO(s.client, zap.NewNop(), "")
}

func (s *CourierTestSuite) ClearCouriersFromElastic(couriersIDs ...string) error {
	for _, id := range couriersIDs {
		if _, err := s.client.Delete().Index(CourierIndex).Id(id).Do(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

//tests
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

func (s *CourierTestSuite) TestGetCourierByIDOK() {
	service := s.GetService()
	name := "Vasya"
	courier := &models.CourierCreate{
		Name: name,
	}
	id := s.CreateCourier(courier)
	res, err := service.GetByID(id)
	if !s.NoError(err) {
		s.Assert().FailNowf("error", "error: %s", err)
	}
	s.Assert().Equal(name, res.Name)
}

func (s *CourierTestSuite) TestGetCourierByIDBadID() {
	service := s.GetService()
	id := "bad id"
	res, err := service.GetByID(id)
	s.Error(err)
	s.Nil(res)
}

func (s *CourierTestSuite) TestUpdateCourierWithoutLocationOK() {
	service := s.GetService()
	name := "Vasya"
	courier := &models.CourierCreate{
		Name: name,
	}
	id := s.CreateCourier(courier)
	phone := "79031234512"
	name = "NewVasya"
	courierUpd := &models.CourierUpdate{
		ID:    &id,
		Name:  &name,
		Phone: &phone,
	}
	res, err := service.Update(courierUpd)
	if !s.Assert().NoError(err) {
		s.FailNow(err.Error())
	}
	s.Assert().IsType(&models.Courier{}, res)
	s.Assert().Equal(name, res.Name)
	s.Assert().Equal(phone, *res.Phone)
}

func (s *CourierTestSuite) TestUpdateCourierWithLocationOK() {
	service := s.GetService()
	name := "Vasya"
	courier := &models.CourierCreate{
		Name: name,
	}
	id := s.CreateCourier(courier)
	phone := "79031234512"
	name = "NewVasya"
	address := "Moscow"
	courierUpd := &models.CourierUpdate{
		ID:    &id,
		Name:  &name,
		Phone: &phone,
		Location: &models.Location{
			GeoPoint: &elastic.GeoPoint{
				Lat:123.023,
				Lon: 123.0123,
			},
			Address: &address,
		},
	}
	res, err := service.Update(courierUpd)
	if !s.Assert().NoError(err) {
		s.FailNow(err.Error())
	}
	s.Assert().Equal(name, res.Name)
	s.Assert().Equal(phone, *res.Phone)
	s.Assert().NotNil(res.Location)
	s.Assert().NotNil(res.LastSeen)
}

func TestIntegrationCouriersDAO(t *testing.T) {
	suite.Run(t, new(CourierTestSuite))
}
