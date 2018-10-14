package interfaces

import "github.com/olivere/elastic"

type GeoRouteInterface interface {
	CreateCourierWithOrder(courierID, orderID string) error
	AddPointToRoute(courierID string, point *elastic.GeoPoint) error
	DeleteCourier(courierID string) error
	GetRoute(courierID, orderID string) ([]*elastic.GeoPoint, error)
}
