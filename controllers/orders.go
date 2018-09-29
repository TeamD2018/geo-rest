package controllers

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
)

func (api *ApiService) GetOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	_, err := uuid.FromString(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.OneOfParameterHaveIncorrectFormat)
		return
	}
	order, err := api.OrdersDAO.Get(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.OneOfParametersNotFound)
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (api *ApiService) UpdateOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	_, err := uuid.FromString(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.OneOfParameterHaveIncorrectFormat)
		return
	}
	var order models.OrderUpdate

	if err := ctx.ShouldBind(&order); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.OneOfParameterHaveIncorrectFormat)
		return
	}
	order.ID = &orderID
	created, err := api.OrdersDAO.Update(&order)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ServerError)
		return
	}
	ctx.JSON(http.StatusOK, created)
}

func (api *ApiService) CreateOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	_, err := uuid.FromString(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.OneOfParameterHaveIncorrectFormat)
		return
	}
	courierID := ctx.Param("courier_id")
	var order models.OrderCreate
	if err := ctx.ShouldBind(&order); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.OneOfParameterHaveIncorrectFormat)
	}
	order.CourierID = &courierID
	created, err := api.OrdersDAO.Create(&order)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ServerError)
		return
	}
	ctx.JSON(http.StatusCreated, created)
}

func (api *ApiService) DeleteOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	_, err := uuid.FromString(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.OneOfParameterHaveIncorrectFormat)
		return
	}
	if err := api.OrdersDAO.Delete(orderID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.EntityNotFound)
		return
	}
	ctx.Status(http.StatusNoContent)
}