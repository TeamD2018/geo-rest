package controllers

import (
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
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
	results, err := api.SuggesterExecutor.Suggest(params.Input)
	if err != nil {
		api.Logger.Error("fail to get suggestions", zap.String("input", input), zap.Error(err))
	}
	ordersRaw := results["orders-engine"]
	couriersRaw := results["couriers-engine"]
	if suggestion, err := models.SuggestionFromRawInput(ordersRaw, couriersRaw);
		err != nil {
		api.Logger.Error("fail to build suggestion from elastic results", zap.Error(err))
	} else {
		ctx.JSON(http.StatusOK, suggestion)
	}
}
