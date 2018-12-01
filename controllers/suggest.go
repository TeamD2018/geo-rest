package controllers

import (
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/TeamD2018/geo-rest/services/suggestions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (api *APIService) Suggest(ctx *gin.Context) {
	var params parameters.GenericSuggestParams
	if err := ctx.BindQuery(&params); err != nil {
		ctx.JSON(models.ErrOneOfParameterHaveIncorrectFormat.HttpStatus(), models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	results, err := api.SuggestionService.Suggest(params.Input)
	if err != nil {
		api.Logger.Error("fail to get suggestions", zap.String("input", params.Input), zap.Error(err))
		return
	}
	orders, _ := results["orders-engine"].([]suggestions.ElasticSuggestResult)
	ordersByPrefix, _ := results["orders-prefix-engine"].([]suggestions.ElasticSuggestResult)
	orders = append(orders, ordersByPrefix...)
	couriersRaw, _ := results["couriers-engine"].([]suggestions.ElasticSuggestResult)
	polygons, _ := results["polygons-engine"].([]*models.OSMPolygonSuggestion)
	suggestion, err := models.SuggestionFromRawInput(orders, couriersRaw, polygons)
	if err != nil {
		api.Logger.Error("fail to build suggestion from elastic results", zap.Error(err))
	}

	if err := api.OrdersCountTracker.Sync(suggestion.Couriers); err != nil {
		api.Logger.Error("fail to sync couriers counters", zap.Error(err))
	}
	ctx.JSON(http.StatusOK, suggestion)
}
