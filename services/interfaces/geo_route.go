package interfaces

import (
	"github.com/TeamD2018/geo-rest/models"
)

type GeoRouteInterface interface {
	CreateCourier(courierID string) error
	AddPointToRoute(courierID string, point *models.PointWithTs) error
	DeleteCourier(courierID string) error
	GetRoute(courierID string) ([]*models.PointWithTs, error)
}
