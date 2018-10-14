package interfaces

import (
	"github.com/TeamD2018/geo-rest/models"
)

type CourierSuggester interface {
	Suggest(field string, prefix string, limit int) (models.Couriers, error)
}
