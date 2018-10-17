package interfaces

import (
	"github.com/TeamD2018/geo-rest/models"
)

type GeoRouteInterface interface {
	CreateCourierWithOrder(courierID, orderID string) error
	AddPointToRoute(courierID string, point *models.PointWithTs) error
	DeleteCourier(courierID string) error
	GetRoute(courierID, orderID string) ([]*models.PointWithTs, error)
}
