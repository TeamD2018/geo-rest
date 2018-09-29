package controllers

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
)

func (api *APIService) CreateCourier(ctx *gin.Context) {
	courier := &models.CourierCreate{}
	if err := ctx.ShouldBind(&courier); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.OneOfParameterHaveIncorrectFormat)
		return
	}
	if res, err := api.CouriersDAO.Create(courier); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ServerError)
		return
	} else {
		ctx.JSON(http.StatusCreated, res)
		return
	}
}

func (api *APIService) GetCourierByID(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	if _, err := uuid.FromString(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.EntityNotFound)
		return
	}
	if courier, err := api.CouriersDAO.GetByID(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.EntityNotFound)
		return
	} else {
		ctx.JSON(http.StatusOK, courier)
		return
	}
}

func (api *APIService) UpdateCourier(ctx *gin.Context) {
	courier := &models.CourierUpdate{}
	if err := ctx.ShouldBind(courier); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.OneOfParameterHaveIncorrectFormat)
		return
	}
	if updated, err := api.CouriersDAO.Update(courier); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ServerError)
		return
	} else {
		ctx.JSON(http.StatusOK, updated)
		return
	}
}

func (api *APIService) DeleteCourier(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	if _, err := uuid.FromString(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.OneOfParameterHaveIncorrectFormat)
		return
	}
	if err := api.CouriersDAO.Delete(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ServerError)
	}
	ctx.Status(http.StatusNoContent)
}
