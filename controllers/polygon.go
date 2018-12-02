package controllers

import (
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api *APIService) GetPolygon(ctx *gin.Context) {
	searchParams := parameters.PolygonQuery{}
	if err := ctx.BindQuery(&searchParams); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	polygon, err := api.RegionResolver.ResolveRegion(searchParams.ToOSMEntity())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	ctx.JSON(http.StatusOK, models.CouriersResponse{Polygon: polygon})
}
