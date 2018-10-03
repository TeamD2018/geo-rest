package interfaces

import "github.com/TeamD2018/geo-rest/models"

type ICouriersDAO interface {
	GetByID(courierID string) (*models.Courier, error)
	GetByName(name string) (models.Couriers, error)
	GetBySquareField(field *models.SquareField) (models.Couriers, error)
	GetByCircleField(field *models.CircleField) (models.Couriers, error)
	Create(courier *models.CourierCreate) (*models.Courier, error)
	Update(courier *models.CourierUpdate) (*models.Courier, error)
	Delete(courierID string) error
}
