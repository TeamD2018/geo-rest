package mocks

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/stretchr/testify/mock"
)

type OrdersCountTrackerMock struct {
	mock.Mock
}

func (octm *OrdersCountTrackerMock) Inc(courierId string) error {
	args := octm.Called(courierId)
	return args.Error(0)
}

func (octm *OrdersCountTrackerMock) Dec(courierId string) error {
	args := octm.Called(courierId)
	return args.Error(0)
}

func (octm *OrdersCountTrackerMock) Sync(couriers models.Couriers) (error) {
	args := octm.Called(couriers)
	for _, courier := range couriers {
		courier.OrdersCount = 0
	}
	return args.Error(0)
}

func (octm *OrdersCountTrackerMock) Drop(courierId string) error {
	args := octm.Called(courierId)
	return args.Error(0)
}
