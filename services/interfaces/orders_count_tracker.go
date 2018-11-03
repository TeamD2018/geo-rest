package interfaces

import "github.com/TeamD2018/geo-rest/models"

type OrdersCountTracker interface {
	Inc(courier_id string) error
	Dec(courier_id string) error
	Sync(ids models.Couriers) (error)
	Drop(courier_id string) error
}
