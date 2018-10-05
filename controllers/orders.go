package controllers

import (
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (api *APIService) GetOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	_, err := uuid.FromString(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	order, err := api.OrdersDAO.Get(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrOneOfParametersNotFound)
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (api *APIService) UpdateOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	_, err := uuid.FromString(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}

	var order models.OrderUpdate
	if err := ctx.ShouldBind(&order); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}

	if err := api.GeoResolver.Resolve(order.Destination);
		err != nil {
		api.Logger.Error("fail to resolve destination",
			zap.String("order_id", orderID),
			zap.Any("dest", order.Destination))
	}
	if err := api.GeoResolver.Resolve(order.Source);
		err != nil {
		api.Logger.Error("fail to resolve source",
			zap.String("order_id", orderID),
			zap.Any("dest", order.Destination))
	}

	order.ID = &orderID
	order.CourierID = nil
	created, err := api.OrdersDAO.Update(&order)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	ctx.JSON(http.StatusOK, created)
}

func (api *APIService) CreateOrder(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	var order models.OrderCreate
	if err := ctx.ShouldBind(&order); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
	}
	order.CourierID = &courierID

	if err := api.GeoResolver.Resolve(&order.Destination);
		err != nil {
		api.Logger.Error("fail to resolve destination",
			zap.String("courier_id", *order.CourierID),
			zap.Any("dest", order.Destination))
	}
	if err := api.GeoResolver.Resolve(&order.Source);
		err != nil {
		api.Logger.Error("fail to resolve source",
			zap.String("courier_id", *order.CourierID),
			zap.Any("dest", order.Destination))
	}
	created, err := api.OrdersDAO.Create(&order)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	ctx.JSON(http.StatusCreated, created)
}

func (api *APIService) AssignNewCourier(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	orderID := ctx.Param("order_id")
	_, err := uuid.FromString(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	updated, err := api.OrdersDAO.Update(&models.OrderUpdate{CourierID: &courierID})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

func (api *APIService) DeleteOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	_, err := uuid.FromString(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	if err := api.OrdersDAO.Delete(orderID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (api *APIService) GetOrdersForCourier(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	_, err := uuid.FromString(courierID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	var params parameters.GetOrdersForCourierParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	orders, err := api.OrdersDAO.GetOrdersForCourier(courierID, params.Since, params.Asc)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}
