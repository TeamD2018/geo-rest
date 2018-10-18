package interfaces

import (
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
)

type CourierSuggester interface {
	Suggest(field string, suggestion *parameters.Suggestion) (models.Couriers, error)
	SetFuzziness(fuzziness int) CourierSuggester
	SetFuzzinessThreshold(threshold int) CourierSuggester
}
