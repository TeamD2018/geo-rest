package interfaces

import "github.com/TeamD2018/geo-rest/models"

type ICouriersDAO interface {
	Get(courierID string) models.Courier
}