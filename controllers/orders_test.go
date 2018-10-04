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
}

func (oc *OrdersControllersTestSuite) SetupSuite() {
	oc.api = &APIService{
		Logger: zap.NewNop(),
	}
	oc.router = gin.Default()
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
			Point: elastic.GeoPointFromLatLon(10, 10),
			Address:  &testSource,
		},
		Destination: models.Location{
			Point: elastic.GeoPointFromLatLon(20, 20),
			Address:  &testDest,
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
			Point: elastic.GeoPointFromLatLon(15, 15),
			Address:  &updatedSource,
		},
	}
}

func TestUnitControllersOrders(t *testing.T) {
	suite.Run(t, new(OrdersControllersTestSuite))
}

func (oc *OrdersControllersTestSuite) BeforeTest(suiteName, testName string) {
	oc.ordersDAOMock = new(mocks.OrdersDAOMock)
}

func (oc *OrdersControllersTestSuite) TestAPIService_CreateOrder_Created() {
	oc.ordersDAOMock.On("Create", mock.Anything).Return(oc.testOrder, nil)
	oc.api.OrdersDAO = oc.ordersDAOMock

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

func (oc *OrdersControllersTestSuite) TestAPIService_UpdateOrder_OK() {
	oc.testOrder.Destination = *oc.testOrderUpdate.Destination
	oc.ordersDAOMock.On("Update", mock.Anything).Return(oc.testOrder, nil)
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
	oc.api.OrdersDAO = oc.ordersDAOMock

	w := httptest.NewRecorder()
	url := fmt.Sprintf("/couriers/%s/orders/%s", oc.testOrder.CourierID, oc.testOrder.ID)
	req, _ := http.NewRequest("DELETE", url, bytes.NewReader([]byte{}))
	oc.router.ServeHTTP(w, req)

	oc.Equal(http.StatusNoContent, w.Code)
}

func toByteReader(source interface{}) *bytes.Reader {
	bin, _ := json.Marshal(source)
	return bytes.NewReader(bin)
}