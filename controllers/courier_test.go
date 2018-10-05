package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CouriersElasticDAOMock struct {
}

func (CouriersElasticDAOMock) GetByID(courierID string) (*models.Courier, error) {
	return &models.Courier{
		ID:   courierID,
		Name: "Vasya",
	}, nil
}

func (CouriersElasticDAOMock) GetByName(name string) (models.Couriers, error) {
	panic("implement me")
}

func (CouriersElasticDAOMock) GetByBoxField(field *models.BoxField) (models.Couriers, error) {
	panic("implement me")
}

func (CouriersElasticDAOMock) GetByCircleField(field *models.CircleField) (models.Couriers, error) {
	panic("implement me")
}

func (mock *CouriersElasticDAOMock) Create(courier *models.CourierCreate) (*models.Courier, error) {
	m := &models.Courier{
		ID:    uuid.NewV4().String(),
		Name:  courier.Name,
		Phone: courier.Phone,
	}
	return m, nil
}

func (CouriersElasticDAOMock) Update(courier *models.CourierUpdate) (*models.Courier, error) {
	panic("implement me")
}

func (CouriersElasticDAOMock) Delete(courierID string) error {
	panic("implement me")
}

type ControllerCouriersTestSuite struct {
	suite.Suite
	api *APIService
	m   *CouriersElasticDAOMock
	r   *gin.Engine
}

func (s *ControllerCouriersTestSuite) setupRouter() {
	g := s.r.Group("/couriers")
	g.POST("/", s.api.CreateCourier)
	g.GET("/:courier_id", s.api.GetCourierByID)
}

func (s *ControllerCouriersTestSuite) SetupSuite() {
	m := &CouriersElasticDAOMock{}
	s.m = m
	s.api = &APIService{
		OrdersDAO:   nil,
		CouriersDAO: m,
	}
	gin.SetMode("test")
	s.r = gin.New()
	s.setupRouter()
}

func (s *ControllerCouriersTestSuite) TestCreateOk() {
	w := httptest.NewRecorder()
	c :=
		`{
			"name": "Vasya",
			"phone": "79031186555"
		}`
	req, err := http.NewRequest(http.MethodPost, "/couriers/", bytes.NewReader([]byte(c)))

	s.Assert().NoError(err)

	s.r.ServeHTTP(w, req)
	res := w.Result()

	s.Assert().Equal(http.StatusCreated, res.StatusCode)
}

func (s *ControllerCouriersTestSuite) TestCreateWithoutName() {
	w := httptest.NewRecorder()
	c :=
		`{
			"phone": "79031186555"
		}`
	req, err := http.NewRequest(http.MethodPost, "/couriers/", bytes.NewReader([]byte(c)))

	s.Assert().NoError(err)

	s.r.ServeHTTP(w, req)
	res := w.Result()
	s.Assert().NoError(err)
	s.Assert().Equalf(http.StatusBadRequest, res.StatusCode, "%s", w.Body.String())
}

func (s *ControllerCouriersTestSuite) TestGetByIDOk() {
	w := httptest.NewRecorder()
	id := uuid.NewV4()
	path := fmt.Sprintf("/couriers/%s", id.String())
	req, err := http.NewRequest(http.MethodGet, path, nil)

	s.Assert().NoError(err)

	s.r.ServeHTTP(w, req)
	res := w.Result()
	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusOK, res.StatusCode)
	var m map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &m)
	s.Assert().NoErrorf(err, "body: %s", w.Body.String())
	s.Assert().Contains(m, "id")
	s.Assert().Equal(id.String(), m["id"])
}

func TestUnitControllersCouriers(t *testing.T) {
	suite.Run(t, new(ControllerCouriersTestSuite))
}
