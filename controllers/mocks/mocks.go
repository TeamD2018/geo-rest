package mocks

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/stretchr/testify/mock"
)

type OrdersDAOMock struct {
	mock.Mock
}

func (o *OrdersDAOMock) Get(orderID string) (*models.Order, error) {
	args := o.Called(orderID)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (o *OrdersDAOMock) Create(order *models.OrderCreate) (*models.Order, error) {
	args := o.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (o *OrdersDAOMock) Update(order *models.OrderUpdate) (*models.Order, error) {
	args := o.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (o *OrdersDAOMock) Delete(orderID string) error {
	args := o.Called(orderID)
	return args.Error(0)
}

func (o *OrdersDAOMock) GetOrdersForCourier(courierID string, since int64, asc bool) (models.Orders, error) {
	args := o.Called(courierID, since, asc)
	return args.Get(0).(models.Orders), args.Error(1)
}

type GeoResolverMock struct {
	mock.Mock
}

func (gr *GeoResolverMock) Resolve(location *models.Location) error {
	args := gr.Called(location)
	return args.Error(0)
}
