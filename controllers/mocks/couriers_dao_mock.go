package mocks

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/stretchr/testify/mock"
)

type CouriersDAOMock struct {
	mock.Mock
}

func (c *CouriersDAOMock) GetByPolygon(polygon models.FlatPolygon, size int, activeOnly bool) (models.Couriers, error) {
	args := c.Called(polygon, size, activeOnly)
	return args.Get(0).(models.Couriers), args.Error(1)
}

func (c *CouriersDAOMock) GetByID(courierID string) (*models.Courier, error) {
	args := c.Called(courierID)
	return args.Get(0).(*models.Courier), args.Error(1)
}

func (c *CouriersDAOMock) GetByName(name string, size int) (models.Couriers, error) {
	args := c.Called(name)
	return args.Get(0).(models.Couriers), args.Error(1)
}

func (c *CouriersDAOMock) GetByBoxField(field *models.BoxField, size int, isActive bool) (models.Couriers, error) {
	args := c.Called(field, size, isActive)
	return args.Get(0).(models.Couriers), args.Error(1)
}

func (c *CouriersDAOMock) GetByCircleField(field *models.CircleField, size int, isActive bool) (models.Couriers, error) {
	args := c.Called(field, size, isActive)
	return args.Get(0).(models.Couriers), args.Error(1)
}

func (c *CouriersDAOMock) Create(courier *models.CourierCreate) (*models.Courier, error) {
	args := c.Called(courier)
	return args.Get(0).(*models.Courier), args.Error(1)
}

func (c *CouriersDAOMock) Update(courier *models.CourierUpdate) (*models.Courier, error) {
	args := c.Called(courier)
	return args.Get(0).(*models.Courier), args.Error(1)
}

func (c *CouriersDAOMock) Delete(courierID string) error {
	args := c.Called(courierID)
	return args.Error(0)
}
func (c *CouriersDAOMock) Exists(courierID string) (bool, error) {
	args := c.Called(courierID)
	return args.Bool(0), args.Error(1)
}
