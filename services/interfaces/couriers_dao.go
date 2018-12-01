package interfaces

import "github.com/TeamD2018/geo-rest/models"

type ICouriersDAO interface {
	GetByID(courierID string) (*models.Courier, error)
	GetByName(name string, size int) (models.Couriers, error)
	GetByBoxField(field *models.BoxField, size int, activeOnly bool) (models.Couriers, error)
	GetByCircleField(field *models.CircleField, size int, activeOnly bool) (models.Couriers, error)
	GetByPolygon(polygon *models.Polygon, size int, activeOnly bool) (models.Couriers, error)
	Create(courier *models.CourierCreate) (*models.Courier, error)
	Update(courier *models.CourierUpdate) (*models.Courier, error)
	Exists(courierID string) (bool, error)
	Delete(courierID string) error
}
