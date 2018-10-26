package models

import (
	"encoding/json"
	"github.com/TeamD2018/geo-rest/services/suggestions"
)

type Suggestion struct {
	Couriers Couriers `json:"couriers"`
	Orders   Orders   `json:"orders"`
}

func SuggestionFromRawInput(ordersRaw, couriersRaw suggestions.EngineSuggestResults) (*Suggestion, error) {
	suggestion := Suggestion{
		Couriers: make(Couriers, 0, len(couriersRaw)),
		Orders:   make(Orders, 0, len(ordersRaw)),
	}
	for _, rawOrder := range ordersRaw {
		var order Order
		if err := json.Unmarshal(*rawOrder.Source, &order); err != nil {
			return nil, err
		}
		suggestion.Orders = append(suggestion.Orders, &order)
	}
	for _, rawCourier := range couriersRaw {
		var courier Courier
		if err := json.Unmarshal(*rawCourier.Source, &courier); err != nil {
			return nil, err
		}
		suggestion.Couriers = append(suggestion.Couriers, &courier)
	}
	return &suggestion, nil
}
