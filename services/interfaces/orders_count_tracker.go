package interfaces

import "github.com/TeamD2018/geo-rest/models"

type OrdersCountTracker interface {
	Inc(courierID string) (error)
	Dec(courierID string) (error)
	DecAndGet(courierID string) (int, error)
	IncAndGet(courierID string) (int, error)
	Sync(ids models.Couriers) (error)
	Drop(courierID string) error
}
