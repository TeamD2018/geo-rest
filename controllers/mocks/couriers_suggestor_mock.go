package mocks

import (
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/stretchr/testify/mock"
)

type CouriersSuggestorMock struct {
	mock.Mock
}

func (os *CouriersSuggestorMock) SetFuzziness(fuzziness int) interfaces.CourierSuggester {
	os.Called(fuzziness)
	return os
}

func (os *CouriersSuggestorMock) SetFuzzinessThreshold(threshold int) interfaces.CourierSuggester {
	os.Called(threshold)
	return os
}

func (os *CouriersSuggestorMock) Suggest(field string, suggestion *parameters.Suggestion) (models.Couriers, error) {
	args := os.Called(field, suggestion)
	return args.Get(0).(models.Couriers), args.Error(1)
}
