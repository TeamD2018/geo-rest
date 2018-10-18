package mocks

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/stretchr/testify/mock"
)

type GeoRouteMock struct {
	mock.Mock
}

func (m *GeoRouteMock) CreateCourier(courierID string) error {
	args := m.Called(courierID)
	return args.Error(0)
}

func (m *GeoRouteMock) AddPointToRoute(courierID string, point *models.PointWithTs) error {
	args := m.Called(courierID, point)
	return args.Error(0)
}

func (m *GeoRouteMock) DeleteCourier(courierID string) error {
	args := m.Called(courierID)
	return args.Error(0)
}

func (m *GeoRouteMock) GetRoute(courierID string) ([]*models.PointWithTs, error) {
	args := m.Called(courierID)
	return args.Get(0).([]*models.PointWithTs), args.Error(1)
}
