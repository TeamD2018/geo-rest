package interfaces

import "github.com/TeamD2018/geo-rest/models"

type OrdersCountTracker interface {
	Inc(courirID string) error
	Dec(courierID string) error
	Sync(ids models.Couriers) (error)
	Drop(courierID string) error
}
