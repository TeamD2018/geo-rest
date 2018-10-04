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

type CouriersDAOMock struct {
	mock.Mock
}

func (c *CouriersDAOMock) GetByID(courierID string) (*models.Courier, error) {
	args := c.Called(courierID)
	return args.Get(0).(*models.Courier), args.Error(1)
}

func (CouriersDAOMock) GetByName(name string) (models.Couriers, error) {
	return nil, nil
}

func (CouriersDAOMock) GetBySquareField(field *models.SquareField) (*models.Courier, error) {
	return nil, nil
}

func (CouriersDAOMock) GetByCircleField(field *models.CircleField) (*models.Courier, error) {
	return nil, nil
}

func (CouriersDAOMock) Create(courier *models.CourierCreate) (*models.Courier, error) {
	return nil, nil
}

func (CouriersDAOMock) Update(courier *models.CourierUpdate) (*models.Courier, error) {
	return nil, nil
}

func (CouriersDAOMock) Delete(courierID string) error {
	return nil
}
