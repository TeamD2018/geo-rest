package controllers

import (
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
)

func (api *APIService) CreateCourier(ctx *gin.Context) {
	var courier models.CourierCreate
	if err := ctx.BindJSON(&courier); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	if res, err := api.CouriersDAO.Create(&courier); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	} else {
		ctx.JSON(http.StatusCreated, res)
		return
	}
}

func (api *APIService) GetCourierByID(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	if _, err := uuid.FromString(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
		return
	}
	if courier, err := api.CouriersDAO.GetByID(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
		return
	} else {
		ctx.JSON(http.StatusOK, courier)
		return
	}
}

func (api *APIService) MiddlewareGeoSearch(ctx *gin.Context) {
	if ctx.Request.URL.Query().Get("r") != "" {
		api.GetCouriersByCircleField(ctx)
		return
	} else {
		return
	}
}

func (api *APIService) UpdateCourier(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	if _, err := uuid.FromString(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	courier := &models.CourierUpdate{}
	if err := ctx.BindJSON(courier); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	courier.ID = &courierID
	if updated, err := api.CouriersDAO.Update(courier); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	} else {
		ctx.JSON(http.StatusOK, updated)
		return
	}
}

func (api *APIService) GetCouriersByCircleField(ctx *gin.Context) {
	circle := parameters.CircleFieldQuery{}
	if err := ctx.BindQuery(&circle); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	if couriers, err := api.CouriersDAO.GetByCircleField(circle.ToCircleField()); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	} else {
		ctx.JSON(http.StatusOK, couriers)
		return
	}
}

func (api *APIService) DeleteCourier(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	if _, err := uuid.FromString(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	if err := api.CouriersDAO.Delete(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
	}
	ctx.Status(http.StatusNoContent)
}
