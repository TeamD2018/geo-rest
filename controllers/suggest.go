package controllers

import (
	"encoding/json"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (api *APIService) Suggest(ctx *gin.Context) {
	input := ctx.Query("input")
	results, err := api.SuggesterExecutor.Suggest(input)
	if err != nil {
		api.Logger.Error("fail to get suggestions", zap.String("input", input), zap.Error(err))
	}
	var suggestion models.Suggestion
	ordersRaw := results["orders-engine"]
	couriersRaw := results["couriers-engine"]
	suggestion.Orders = make(models.Orders, 0, len(ordersRaw))
	suggestion.Couriers = make(models.Couriers, 0, len(couriersRaw))
	for _, rawOrder := range ordersRaw {
		var order models.Order
		if err := json.Unmarshal(*rawOrder.Source, &order); err != nil {
			ctx.JSON(models.ErrUnmarshalJSON.HttpCode, models.ErrUnmarshalJSON)
			return
		}
		suggestion.Orders = append(suggestion.Orders, &order)
	}
	for _, rawCourier := range couriersRaw {
		var courier models.Courier
		if err := json.Unmarshal(*rawCourier.Source, &courier); err != nil {
			ctx.JSON(models.ErrUnmarshalJSON.HttpCode, models.ErrUnmarshalJSON)
			return
		}
		suggestion.Couriers = append(suggestion.Couriers, &courier)
	}
	ctx.JSON(http.StatusOK, suggestion)
}
