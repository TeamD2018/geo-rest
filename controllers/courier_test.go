package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TeamD2018/geo-rest/controllers/mocks"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type ControllerCouriersTestSuite struct {
	suite.Suite
	api               *APIService
	router            *gin.Engine
	testCourier       *models.Courier
	testCourierCreate *models.CourierCreate
	testCourierUpdate *models.CourierUpdate
	geoRouteMock      *mocks.GeoRouteMock
	couriersDAOMock   *mocks.CouriersDAOMock
}

func (ts *ControllerCouriersTestSuite) SetupSuite() {
	geoResolverMock := new(mocks.GeoResolverMock)
	geoResolverMock.On("Resolve", mock.Anything).Return(nil)
	ts.api = &APIService{
		Logger:      zap.NewNop(),
		GeoResolver: geoResolverMock,
	}
	gin.DisableConsoleColor()
	ts.router = gin.New()
	SetupRouters(ts.router, ts.api)
	testPhone := "Test phone"

	ts.testCourier = &models.Courier{
		ID:    "550e8400-e29b-41d4-a716-446655440000",
		Name:  "Test Name",
		Phone: &testPhone,
	}

	ts.testCourierCreate = &models.CourierCreate{
		Name:  ts.testCourier.Name,
		Phone: ts.testCourier.Phone,
	}

	testLocation := "Test location"

	ts.testCourierUpdate = &models.CourierUpdate{
		Location: &models.Location{
			Point:   elastic.GeoPointFromLatLon(10, 10),
			Address: &testLocation,
		},
	}
}

func TestUnitControllersCouriers(t *testing.T) {
	suite.Run(t, new(ControllerCouriersTestSuite))
}

func (ts *ControllerCouriersTestSuite) BeforeTest(suiteName, testName string) {
	ts.couriersDAOMock = new(mocks.CouriersDAOMock)
	ts.geoRouteMock = new(mocks.GeoRouteMock)
}

func (ts *ControllerCouriersTestSuite) TestAPIService_CreateCourier_Created() {
	ts.couriersDAOMock.On("Create", mock.Anything).Return(ts.testCourier, nil)
	ts.api.CouriersDAO = ts.couriersDAOMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers")
	req, _ := http.NewRequest("POST", url, toByteReader(ts.testCourierCreate))
	ts.router.ServeHTTP(w, req)

	var got models.Courier
	err := json.Unmarshal(w.Body.Bytes(), &got)

	ts.NoError(err)
	ts.Equal(201, w.Code)
	ts.Equal(ts.testCourier, &got)
}

func (ts *ControllerCouriersTestSuite) TestAPIService_GetCourierById_OK() {
	ts.couriersDAOMock.On("GetByID", ts.testCourier.ID).Return(ts.testCourier, nil)
	ts.api.CouriersDAO = ts.couriersDAOMock

	url := fmt.Sprintf("/couriers/%s", ts.testCourier.ID)
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	var got models.Courier
	err := json.Unmarshal(w.Body.Bytes(), &got)

	ts.NoError(err)
	ts.Equal(http.StatusOK, w.Code)
	ts.Equal(ts.testCourier, &got)
}

func (ts *ControllerCouriersTestSuite) TestAPIService_UpdateCourier_OK() {
	ts.testCourier.Location = ts.testCourierUpdate.Location
	ts.couriersDAOMock.On("Update", mock.Anything).Return(ts.testCourier, nil)
	ts.geoRouteMock.On("AddPointToRoute", mock.Anything, mock.Anything).Return(nil)
	ts.api.CouriersDAO = ts.couriersDAOMock
	ts.api.CourierRouteDAO = ts.geoRouteMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s", ts.testCourier.ID)
	req, _ := http.NewRequest("PUT", url, toByteReader(ts.testCourierUpdate))
	ts.router.ServeHTTP(w, req)

	var got models.Courier
	err := json.Unmarshal(w.Body.Bytes(), &got)

	ts.NoError(err)
	ts.Equal(http.StatusOK, w.Code)
	ts.Equal(ts.testCourier, &got)
}

func (ts *ControllerCouriersTestSuite) TestAPIService_DeleteOrder_NoContent() {
	ts.couriersDAOMock.On("Delete", ts.testCourier.ID).Return(nil)
	ts.api.CouriersDAO = ts.couriersDAOMock
	ts.geoRouteMock.On("DeleteCourier", ts.testCourier.ID).Return(nil)
	ts.api.CourierRouteDAO = ts.geoRouteMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s", ts.testCourier.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	ts.router.ServeHTTP(w, req)

	ts.Equal(http.StatusNoContent, w.Code)
}

func (ts *ControllerCouriersTestSuite) TestAPIService_GetCouriersByCircleField() {
	circleField := &models.CircleField{
		Center: elastic.GeoPointFromLatLon(10, 10),
		Radius: 10,
	}

	testCouriers := append(models.Couriers{}, ts.testCourier)

	ts.couriersDAOMock.On("GetByCircleField", circleField).Return(testCouriers, nil)
	ts.api.CouriersDAO = ts.couriersDAOMock

	v := url.Values{}
	v.Add("radius", "10")
	v.Add("lat", "10")
	v.Add("lon", "10")

	uri := fmt.Sprintf("/couriers?%s", v.Encode())
	req, _ := http.NewRequest("GET", uri, bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	var got models.Couriers
	err := json.Unmarshal(w.Body.Bytes(), &got)

	ts.NoError(err)
	ts.Equal(http.StatusOK, w.Code)
	ts.Equal(testCouriers, got)
}

func (ts *ControllerCouriersTestSuite) TestAPIService_GetCouriersByBoxField() {
	boxField := &models.BoxField{
		TopLeftPoint:     elastic.GeoPointFromLatLon(10, 10),
		BottomRightPoint: elastic.GeoPointFromLatLon(20, 20),
	}

	testCouriers := append(models.Couriers{}, ts.testCourier)

	ts.couriersDAOMock.On("GetByBoxField", boxField).Return(testCouriers, nil)
	ts.api.CouriersDAO = ts.couriersDAOMock

	v := url.Values{}
	v.Add("top_left_lat", "10")
	v.Add("top_left_lon", "10")
	v.Add("bottom_right_lat", "20")
	v.Add("bottom_right_lon", "20")

	uri := fmt.Sprintf("/couriers?%s", v.Encode())
	req, _ := http.NewRequest("GET", uri, bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	var got models.Couriers
	err := json.Unmarshal(w.Body.Bytes(), &got)

	ts.NoError(err)
	ts.Equal(http.StatusOK, w.Code)
	ts.Equal(testCouriers, got)
}
