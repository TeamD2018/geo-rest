package mocks

import (
	"context"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/stretchr/testify/mock"
)

type OrdersDAOMock struct {
	mock.Mock
}

func (o *OrdersDAOMock) Get(orderID string) (*models.Order, error) {
	args := o.Called(orderID)
	v := args.Get(0)
	err := args.Error(1)
	switch v.(type) {
	case *models.Order:
		if v == nil {
			return nil, err
		}
		return v.(*models.Order), err
	default:
		return nil, err
	}
}

func (o *OrdersDAOMock) Create(order *models.OrderCreate) (*models.Order, error) {
	args := o.Called(order)
	v := args.Get(0)
	err := args.Error(1)
	switch v.(type) {
	case *models.Order:
		if v == nil {
			return nil, err
		}
		return v.(*models.Order), err
	default:
		return nil, err
	}
}

func (o *OrdersDAOMock) Update(order *models.OrderUpdate) (*models.Order, error) {
	args := o.Called(order)
	v := args.Get(0)
	err := args.Error(1)
	switch v.(type) {
	case *models.Order:
		if v == nil {
			return nil, err
		}
		return v.(*models.Order), err
	default:
		return nil, err
	}
}

func (o *OrdersDAOMock) Delete(orderID string) error {
	args := o.Called(orderID)
	return args.Error(0)
}

func (o *OrdersDAOMock) GetOrdersForCourier(courierID string, since int64, asc bool) (models.Orders, error) {
	args := o.Called(courierID, since, asc)
	v := args.Get(0)
	err := args.Error(1)
	switch v.(type) {
	case models.Orders:
		if v == nil {
			return nil, err
		}
		return v.(models.Orders), err
	default:
		return nil, err
	}
}

type GeoResolverMock struct {
	mock.Mock
}

func (gr *GeoResolverMock) Resolve(location *models.Location, ctx context.Context) error {
	args := gr.Called(location,ctx)
	return args.Error(0)
}
