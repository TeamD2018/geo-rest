package interfaces

import (
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
)

type IOrdersDao interface {
	Get(orderID string) (*models.Order, error)
	Create(order *models.OrderCreate) (*models.Order, error)
	Update(order *models.OrderUpdate) (*models.Order, error)
	Delete(orderID string) error
	GetOrdersForCourier(courierID string,
		since int64,
		asc parameters.DirectionFlag,
		excludeDelivered parameters.DeliveredFlag) (models.Orders, error)
}
