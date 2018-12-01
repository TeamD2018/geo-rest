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
		"TestGetCourierByIDNoExistID",
		"TestUpdateCourierWithoutLocationOK",
		"TestUpdateCourierWithLocationOK",
		"TestUpdateCourierNoExistID",
		"TestDeleteCourierOK",
		"TestDeleteCourierNoExistID",
		"TestGetCourierByID",
		"TestGetCouriersByCircleFieldOK",
		"TestGetCouriersByCircleFieldEmpty",
		"TestGetCouriersByBoxFieldOK",
		"TestGetCouriersByBoxFieldEmpty",
		"TestGetCouriersByPolygonOK",
		"TestExistsCourierNotFound",
		"TestExistsCourierOK",
		"TestSuggestByPhoneOK",
		"TestSuggestByNameOK",
		"TestSuggestByPhoneFuzzyOK",
		"TestGetCouriersByBoxFieldActiveOnly",
	}
	testsWithDeleteIndex = []string{
		"TestCreateCourierWithNameAndPhone",
		"TestCreateCourierWithName",
		"TestGetCourierByIDOK",
		"TestGetCourierByIDNoExistID",
		"TestUpdateCourierWithoutLocationOK",
		"TestUpdateCourierWithLocationOK",
		"TestUpdateCourierNoExistID",
		"TestDeleteCourierOK",
		"TestDeleteCourierNoExistID",
		"TestCouriersElasticDAO_EnsureMapping",
		"TestGetCouriersByCircleFieldOK",
		"TestGetCouriersByCircleFieldEmpty",
		"TestGetCouriersByBoxFieldOK",
		"TestGetCouriersByBoxFieldEmpty",
		"TestGetCouriersByPolygonOK",
		"TestExistsCourierNotFound",
		"TestExistsCourierOK",
		"TestSuggestByPhoneOK",
		"TestSuggestByNameOK",
		"TestSuggestByPhoneFuzzyOK",
		"TestGetCouriersByBoxFieldActiveOnly",
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

func (s *CourierTestSuite) UpdateCourier(courier *models.CourierUpdate) {
	service := s.GetService()
	_, err := service.Update(courier)
	if !s.Assert().NoError(err) {
		s.Assert().FailNow(err.Error())
	}
	return
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

	resource, err := pool.Run("docker.elastic.co/elasticsearch/elasticsearch", "6.3.2", []string{"discovery.type=single-node"})
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
	return NewCouriersElasticDAO(s.client, zap.NewNop(), "", DefaultCouriersReturnSize)
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

func (s *CourierTestSuite) TestGetCourierByIDNoExistID() {
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

func (s *CourierTestSuite) TestUpdateCourierNoExistID() {
	service := s.GetService()
	name := "Vasya"
	courier := &models.CourierCreate{
		Name: name,
	}
	s.CreateCourier(courier)
	phone := "79031234512"
	name = "NewVasya"
	id := "NoExistID"
	courierUpd := &models.CourierUpdate{
		ID:    &id,
		Name:  &name,
		Phone: &phone,
	}
	res, err := service.Update(courierUpd)
	if !s.Assert().Error(err) {
		s.FailNow(err.Error())
	}
	s.Assert().Nil(res)
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
			Point: &elastic.GeoPoint{
				Lat: 70.0123,
				Lon: 70.0123,
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

func (s *CourierTestSuite) TestDeleteCourierOK() {
	service := s.GetService()
	name := "Vasya"
	courier := &models.CourierCreate{
		Name: name,
	}
	id := s.CreateCourier(courier)
	err := service.Delete(id)
	if !s.NoError(err) {
		s.Assert().FailNowf("error", "error: %s", err)
	}
}

func (s *CourierTestSuite) TestDeleteCourierNoExistID() {
	service := s.GetService()
	id := "NoExistID"
	err := service.Delete(id)
	if !s.Assert().Error(err) {
		s.Assert().FailNowf("error", "error: %s", err)
	}
}

func (s *CourierTestSuite) TestGetCouriersByCircleFieldOK() {
	service := s.GetService()
	name := "Vasya"
	phone := "79031189023"
	courier := &models.CourierCreate{
		Name:  name,
		Phone: &phone,
	}
	id := s.CreateCourier(courier)
	courierUpd := &models.CourierUpdate{
		ID: &id,
		Location: &models.Location{
			Point: elastic.GeoPointFromLatLon(70.0, 70.0),
		},
	}
	s.UpdateCourier(courierUpd)
	s.client.Refresh(CourierIndex).Do(context.Background())
	res, err := service.GetByCircleField(&models.CircleField{
		Center: elastic.GeoPointFromLatLon(70.00005, 70.00005),
		Radius: 1000,
	}, 0, false)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(res)
	s.Assert().Len(res, 1)
}

func (s *CourierTestSuite) TestGetCouriersByCircleFieldEmpty() {
	service := s.GetService()
	name := "Vasya"
	phone := "79031189023"
	courier := &models.CourierCreate{
		Name:  name,
		Phone: &phone,
	}
	id := s.CreateCourier(courier)
	courierUpd := &models.CourierUpdate{
		ID: &id,
		Location: &models.Location{
			Point: elastic.GeoPointFromLatLon(70.0, 70.0),
		},
	}
	s.UpdateCourier(courierUpd)
	s.client.Refresh(CourierIndex).Do(context.Background())
	res, err := service.GetByCircleField(&models.CircleField{
		Center: elastic.GeoPointFromLatLon(1, 1),
		Radius: 1,
	}, 0, false)
	s.Assert().NoError(err)
	s.Assert().Empty(res)
}

func (s *CourierTestSuite) TestGetCouriersByBoxFieldOK() {
	service := s.GetService()
	name := "Vasya"
	phone := "79031189023"
	courier := &models.CourierCreate{
		Name:  name,
		Phone: &phone,
	}
	id := s.CreateCourier(courier)
	courierUpd := &models.CourierUpdate{
		ID: &id,
		Location: &models.Location{
			Point: elastic.GeoPointFromLatLon(70.0, 70.0),
		},
	}
	s.UpdateCourier(courierUpd)
	s.client.Refresh(CourierIndex).Do(context.Background())
	res, err := service.GetByBoxField(&models.BoxField{
		TopLeftPoint:     elastic.GeoPointFromLatLon(71.0, 69.0),
		BottomRightPoint: elastic.GeoPointFromLatLon(0, 0),
	}, 0, false)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(res)
	s.Assert().Len(res, 1)
}

func (s *CourierTestSuite) TestGetCouriersByBoxFieldEmpty() {
	service := s.GetService()
	name := "Vasya"
	phone := "79031189023"
	courier := &models.CourierCreate{
		Name:  name,
		Phone: &phone,
	}
	id := s.CreateCourier(courier)
	courierUpd := &models.CourierUpdate{
		ID: &id,
		Location: &models.Location{
			Point: elastic.GeoPointFromLatLon(70.0, 70.0),
		},
	}
	s.UpdateCourier(courierUpd)
	s.client.Refresh(CourierIndex).Do(context.Background())
	res, err := service.GetByBoxField(&models.BoxField{
		TopLeftPoint:     elastic.GeoPointFromLatLon(1, 1),
		BottomRightPoint: elastic.GeoPointFromLatLon(0, 0),
	}, 0, false)
	s.Assert().NoError(err)
	s.Assert().Empty(res)
}

func (s *CourierTestSuite) TestGetCouriersByBoxFieldActiveOnly() {
	service := s.GetService()
	phone := "79031189023"
	activeCourier := &models.CourierCreate{
		Name:     "Active",
		Phone:    &phone,
		IsActive: true,
	}
	inactiveCourier := &models.CourierCreate{
		Name:     "Inactive",
		Phone:    &phone,
		IsActive: false,
	}
	idActive := s.CreateCourier(activeCourier)
	idInactive := s.CreateCourier(inactiveCourier)

	courierUpd := &models.CourierUpdate{
		ID: &idActive,
		Location: &models.Location{
			Point: elastic.GeoPointFromLatLon(70.0, 70.0),
		},
	}
	s.UpdateCourier(courierUpd)
	courierUpd.ID = &idInactive
	s.UpdateCourier(courierUpd)
	s.client.Refresh(CourierIndex).Do(context.Background())
	res, err := service.GetByBoxField(&models.BoxField{
		TopLeftPoint:     elastic.GeoPointFromLatLon(71.0, 69.0),
		BottomRightPoint: elastic.GeoPointFromLatLon(0, 0),
	}, 0, true)
	if !s.NoError(err) {
		return
	}
	if !s.Len(res, 1) {
		return
	}
	s.Equal(res[0].ID, idActive)
}

func (s *CourierTestSuite) TestGetCouriersByPolygonOK() {
	service := s.GetService()
	name := "Vasya"
	phone := "79031189023"
	courier := &models.CourierCreate{
		Name:  name,
		Phone: &phone,
	}
	id := s.CreateCourier(courier)
	courierUpd := &models.CourierUpdate{
		ID: &id,
		Location: &models.Location{
			Point: elastic.GeoPointFromLatLon(55.763615, 37.607687),
		},
	}
	s.UpdateCourier(courierUpd)
	s.client.Refresh(CourierIndex).Do(context.Background())
	polygon := &models.Polygon{
		Points: []*elastic.GeoPoint{
			elastic.GeoPointFromLatLon(56.514792, 36.375407),
			elastic.GeoPointFromLatLon(56.673754, 39.858289),
			elastic.GeoPointFromLatLon(54.692269, 38.979383),
			elastic.GeoPointFromLatLon(55.146880, 36.122938),
			elastic.GeoPointFromLatLon(56.514792, 36.375407),
		},
	}
	res, err := service.GetByPolygon(polygon, 1, false)
	if !s.NoError(err) || !s.NotEmpty(res) {
		return
	}
}

func (s *CourierTestSuite) TestExistsCourierOK() {
	service := s.GetService()
	name := "Vasya"
	phone := "79031189023"
	courier := &models.CourierCreate{
		Name:  name,
		Phone: &phone,
	}
	id := s.CreateCourier(courier)
	s.client.Refresh(service.index).Do(context.Background())
	isExists, err := service.Exists(id)
	s.NoError(err)
	s.True(isExists)
}

func (s *CourierTestSuite) TestExistsCourierNotFound() {
	service := s.GetService()
	isExists, err := service.Exists(uuid.NewV4().String())
	s.NoError(err)
	s.False(isExists)
}

func TestIntegrationCouriersDAO(t *testing.T) {
	suite.Run(t, new(CourierTestSuite))
}
