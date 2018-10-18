package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeamD2018/geo-rest/controllers/mocks"
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

type OrdersControllersTestSuite struct {
	suite.Suite
	api             *APIService
	router          *gin.Engine
	testCourier     *models.Courier
	testOrder       *models.Order
	testOrderCreate *models.OrderCreate
	testOrderUpdate *models.OrderUpdate
	ordersDAOMock   *mocks.OrdersDAOMock
	geoRouteMock    *mocks.GeoRouteMock
	suggestorMock   *mocks.CouriersSuggestorMock
}

func (oc *OrdersControllersTestSuite) SetupSuite() {
	oc.api = &APIService{
		Logger: zap.NewNop(),
	}
	gin.DisableConsoleColor()
	oc.router = gin.New()
	SetupRouters(oc.router, oc.api)
	testPhone := "Test phone"
	var lastSeen int64 = 100

	oc.testCourier = &models.Courier{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Name:     "Test Name",
		Phone:    &testPhone,
		LastSeen: &lastSeen,
	}

	testSource := "Test source"
	testDest := "Test dest"
	oc.testOrder = &models.Order{
		ID:        "660e8400-e29b-41d4-a716-446655440000",
		CourierID: "550e8400-e29b-41d4-a716-446655440000",
		Source: models.Location{
			Point:   elastic.GeoPointFromLatLon(10, 10),
			Address: &testSource,
		},
		Destination: models.Location{
			Point:   elastic.GeoPointFromLatLon(20, 20),
			Address: &testDest,
		},
		CreatedAt: 100,
	}
	oc.testOrderCreate = &models.OrderCreate{
		CourierID:   &oc.testOrder.CourierID,
		Destination: oc.testOrder.Destination,
		Source:      oc.testOrder.Source,
	}
	updatedSource := "updated source"
	oc.testOrderUpdate = &models.OrderUpdate{
		Destination: &models.Location{
			Point:   elastic.GeoPointFromLatLon(15, 15),
			Address: &updatedSource,
		},
	}
}

func TestUnitControllersOrders(t *testing.T) {
	suite.Run(t, new(OrdersControllersTestSuite))
}

func (oc *OrdersControllersTestSuite) BeforeTest(suiteName, testName string) {
	oc.ordersDAOMock = new(mocks.OrdersDAOMock)
	geoResolverMock := new(mocks.GeoResolverMock)
	geoResolverMock.On("Resolve", mock.Anything, mock.Anything).Return(nil)
	geoResolverMock.On("GetOrdersForCourier", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.AnythingOfType("models.Orders"), mock.AnythingOfType("error"))
	oc.api.GeoResolver = geoResolverMock
	oc.suggestorMock = new(mocks.CouriersSuggestorMock)
	oc.geoRouteMock = new(mocks.GeoRouteMock)
}

func (oc *OrdersControllersTestSuite) TestAPIService_CreateOrder_Created() {
	oc.ordersDAOMock.On("Create", mock.Anything).Return(oc.testOrder, nil)
	oc.geoRouteMock.On("CreateCourier", mock.Anything).Return(nil)
	oc.api.OrdersDAO = oc.ordersDAOMock
	oc.api.CourierRouteDAO = oc.geoRouteMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders", oc.testOrder.CourierID)
	req, _ := http.NewRequest("POST", url, toByteReader(oc.testOrderCreate))
	oc.router.ServeHTTP(w, req)

	var got models.Order
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(201, w.Code)
	oc.Equal(oc.testOrder, &got)
}

func (oc *OrdersControllersTestSuite) TestAPIService_CreateOrder_Created_If_Resolver_Failed() {
	oc.ordersDAOMock.On("Create", mock.Anything).Return(oc.testOrder, nil)
	oc.api.OrdersDAO = oc.ordersDAOMock
	georesolver := new(mocks.GeoResolverMock)
	georesolver.On("Resolve", mock.Anything, mock.Anything).Return(errors.New("test error"))
	georesolver.On("GetOrdersForCourier", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.AnythingOfType("models.Orders"), mock.AnythingOfType("error"))
	oc.api.GeoResolver = georesolver

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders", oc.testOrder.CourierID)
	req, _ := http.NewRequest("POST", url, toByteReader(oc.testOrderCreate))
	oc.router.ServeHTTP(w, req)

	var got models.Order
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(201, w.Code)
	oc.Equal(oc.testOrder, &got)
}
func (oc *OrdersControllersTestSuite) TestAPIService_CreateOrder_EntityNotFound() {
	oc.ordersDAOMock.On("Create", mock.Anything).Return(nil, &models.ErrEntityNotFound)
	oc.api.OrdersDAO = oc.ordersDAOMock
	georesolver := new(mocks.GeoResolverMock)
	georesolver.On("Resolve", mock.Anything, mock.Anything).Return(nil)
	oc.api.GeoResolver = georesolver

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders", oc.testOrder.CourierID)
	req, _ := http.NewRequest("POST", url, toByteReader(oc.testOrderCreate))
	oc.router.ServeHTTP(w, req)

	var got models.Error
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(models.ErrEntityNotFound.HttpStatus(), w.Code)
	oc.Equal(models.ErrEntityNotFound.Code, got.Code)
}

func (oc *OrdersControllersTestSuite) TestAPIService_CreateOrder_UnexpectedError() {
	oc.ordersDAOMock.On("Create", mock.Anything).Return(nil, errors.New("unexpected"))
	oc.api.OrdersDAO = oc.ordersDAOMock
	georesolver := new(mocks.GeoResolverMock)
	georesolver.On("Resolve", mock.Anything, mock.Anything).Return(nil)
	oc.api.GeoResolver = georesolver

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders", oc.testOrder.CourierID)
	req, _ := http.NewRequest("POST", url, toByteReader(oc.testOrderCreate))
	oc.router.ServeHTTP(w, req)

	var got models.Error
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(models.ErrServerError.HttpStatus(), w.Code)
	oc.Equal(models.ErrServerError.Code, got.Code)
}

func (oc *OrdersControllersTestSuite) TestAPIService_GetOrder_OK() {
	oc.ordersDAOMock.On("Get", oc.testOrder.ID).Return(oc.testOrder, nil)
	oc.api.OrdersDAO = oc.ordersDAOMock

	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	oc.router.ServeHTTP(w, req)

	var got models.Order
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(http.StatusOK, w.Code)
	oc.Equal(oc.testOrder, &got)
}

func (oc *OrdersControllersTestSuite) TestAPIService_GetOrder_NotFound() {
	oc.ordersDAOMock.On("Get", oc.testOrder.ID).Return(oc.testOrder, &models.ErrEntityNotFound)
	oc.api.OrdersDAO = oc.ordersDAOMock

	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	oc.router.ServeHTTP(w, req)

	var got models.Error
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(models.ErrEntityNotFound.HttpStatus(), w.Code)
	oc.Equal(models.ErrEntityNotFound.Code, got.Code)
}

func (oc *OrdersControllersTestSuite) TestAPIService_GetOrder_UnexpectedError() {
	oc.ordersDAOMock.On("Get", oc.testOrder.ID).Return(oc.testOrder, errors.New("unexpected"))
	oc.api.OrdersDAO = oc.ordersDAOMock

	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("GET", url, bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()
	oc.router.ServeHTTP(w, req)

	var got models.Error
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(models.ErrServerError.HttpStatus(), w.Code)
	oc.Equal(models.ErrServerError.Code, got.Code)
}

func (oc *OrdersControllersTestSuite) TestAPIService_UpdateOrder_NotFound() {
	oc.testOrder.Destination = *oc.testOrderUpdate.Destination
	oc.ordersDAOMock.On("Update", mock.Anything).Return(oc.testOrder, &models.ErrEntityNotFound)
	oc.api.OrdersDAO = oc.ordersDAOMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("PUT", url, toByteReader(oc.testOrderUpdate))
	oc.router.ServeHTTP(w, req)

	var got models.Error
	err := json.Unmarshal(w.Body.Bytes(), &got)
	oc.NoError(err)
	oc.Equal(models.ErrEntityNotFound.HttpStatus(), w.Code)
	oc.Equal(models.ErrEntityNotFound.Code, got.Code)
}

func (oc *OrdersControllersTestSuite) TestAPIService_UpdateOrder_UnexpectedError() {
	oc.testOrder.Destination = *oc.testOrderUpdate.Destination
	oc.ordersDAOMock.On("Update", mock.Anything).Return(oc.testOrder, errors.New("unexpected"))
	oc.api.OrdersDAO = oc.ordersDAOMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("PUT", url, toByteReader(oc.testOrderUpdate))
	oc.router.ServeHTTP(w, req)

	var got models.Error
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(models.ErrServerError.HttpStatus(), w.Code)
	oc.Equal(models.ErrServerError.Code, got.Code)
}

func (oc *OrdersControllersTestSuite) TestAPIService_UpdateOrder_OK_If_Resolver_Failed() {
	oc.testOrder.Destination = *oc.testOrderUpdate.Destination
	oc.ordersDAOMock.On("Update", mock.Anything).Return(oc.testOrder, nil)
	oc.ordersDAOMock.On("GetOrdersForCourier", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(models.Orders{oc.testOrder}, nil)
	oc.api.OrdersDAO = oc.ordersDAOMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("PUT", url, toByteReader(oc.testOrderUpdate))
	oc.router.ServeHTTP(w, req)

	var got models.Order
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(http.StatusOK, w.Code)
	oc.Equal(oc.testOrder, &got)
}

func (oc *OrdersControllersTestSuite) TestAPIService_DeleteOrder_NoContent() {

	oc.ordersDAOMock.On("Delete", oc.testOrder.ID).Return(nil)
	oc.geoRouteMock.On("DeleteCourier", oc.testOrder.CourierID).Return(nil)
	oc.api.OrdersDAO = oc.ordersDAOMock
	oc.api.CourierRouteDAO = oc.geoRouteMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	oc.router.ServeHTTP(w, req)

	oc.Equal(http.StatusNoContent, w.Code)
}

func (oc *OrdersControllersTestSuite) TestAPIService_DeleteOrder_NotFound() {

	oc.ordersDAOMock.On("Delete", oc.testOrder.ID).Return(&models.ErrEntityNotFound)
	oc.api.OrdersDAO = oc.ordersDAOMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	oc.router.ServeHTTP(w, req)

	var got models.Error
	err := json.Unmarshal(w.Body.Bytes(), &got)
	oc.NoError(err)
	oc.Equal(models.ErrEntityNotFound.HttpStatus(), w.Code)
	oc.Equal(models.ErrEntityNotFound.Code, got.Code)
}

func (oc *OrdersControllersTestSuite) TestAPIService_DeleteOrder_UnexpectedError() {

	oc.ordersDAOMock.On("Delete", oc.testOrder.ID).Return(errors.New("unexpected"))
	oc.api.OrdersDAO = oc.ordersDAOMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	oc.router.ServeHTTP(w, req)

	var got models.Error
	err := json.Unmarshal(w.Body.Bytes(), &got)
	oc.NoError(err)
	oc.Equal(models.ErrServerError.HttpStatus(), w.Code)
	oc.Equal(models.ErrServerError.Code, got.Code)
}

func toByteReader(source interface{}) *bytes.Reader {
	bin, _ := json.Marshal(source)
	return bytes.NewReader(bin)
}

func (oc *OrdersControllersTestSuite) TestAPIService_GetOrdersForCourier_OK() {
	oc.ordersDAOMock.On("GetOrdersForCourier", oc.testOrder.CourierID, int64(0), parameters.WithUpperThreshold, parameters.IncludeDelivered).Return(models.Orders{oc.testOrder}, nil)
	oc.api.OrdersDAO = oc.ordersDAOMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders", oc.testOrder.CourierID)
	req, _ := http.NewRequest("GET", url, nil)
	oc.router.ServeHTTP(w, req)

	var got models.Orders
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(http.StatusOK, w.Code)
	oc.Contains(got, oc.testOrder)
}

func (oc *OrdersControllersTestSuite) TestAPIService_SuggestCouriers_OK() {
	oc.suggestorMock.On("Suggest", mock.Anything, mock.Anything).Return(models.Couriers{oc.testCourier}, nil)
	oc.api.CourierSuggester = oc.suggestorMock

	params := parameters.Suggestion{Prefix: "Test", Limit: 200}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/suggestions/couriers", toByteReader(params))
	oc.router.ServeHTTP(w, req)

	var got models.Couriers
	err := json.Unmarshal(w.Body.Bytes(), &got)

	oc.NoError(err)
	oc.Equal(http.StatusOK, w.Code)
	oc.Contains(got, oc.testCourier)
}
