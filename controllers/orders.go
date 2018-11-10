package controllers

import (
	"context"
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
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
		api.Logger.Error("fail to get order",
			zap.String("order_id", orderID),
			zap.Error(err))
		switch err.(type) {
		case *elastic.Error:
			err := err.(*elastic.Error)
			if err.Status == 404 {
				ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
				return
			}
		case *models.Error:
			err := err.(*models.Error)
			ctx.AbortWithStatusJSON(err.HttpStatus(), err)
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (api *APIService) UpdateOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	_, err := uuid.FromString(orderID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat.SetParameter(orderID))
		return
	}

	courierID := ctx.Param("courier_id")
	_, err = uuid.FromString(courierID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParametersNotFound)
		return
	}
	var order models.OrderUpdate
	if err := ctx.ShouldBind(&order); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	inCtx := context.Background()
	api.Logger.Debug("update data", zap.Any("orderUpdate", order),
		zap.String("order_id", orderID),
		zap.String("courier_id", courierID))
	if err := api.GeoResolver.Resolve(order.Destination, inCtx); err != nil {
		api.Logger.Error("fail to resolve destination",
			zap.Error(err),
			zap.String("order_id", orderID),
			zap.Any("dest", order.Destination))
	}
	if err := api.GeoResolver.Resolve(order.Source, inCtx); err != nil {
		api.Logger.Error("fail to resolve source",
			zap.Error(err),
			zap.String("order_id", orderID),
			zap.Any("dest", order.Destination))
	}

	order.ID = &orderID
	order.CourierID = nil
	created, err := api.OrdersDAO.Update(&order)
	if err != nil {
		api.Logger.Error("fail to update order",
			zap.String("courier_id", courierID),
			zap.String("order_id", orderID),
			zap.Error(err))
		switch err.(type) {
		case *elastic.Error:
			err := err.(*elastic.Error)
			if err.Status == 404 {
				ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
				return
			}
		case *models.Error:
			err := err.(*models.Error)
			ctx.AbortWithStatusJSON(err.HttpStatus(), err)
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	if order.DeliveredAt != nil {
		ordersCount, err := api.OrdersCountTracker.DecAndGet(courierID)
		if err == nil && ordersCount == 0 {
			if err := api.CourierRouteDAO.DeleteCourier(courierID); err != nil {
				api.Logger.Error("fail to cleanup courier route", zap.Error(err))
			}
		}
		if err != nil {
			api.Logger.Error("fail to decrement courier orders count", zap.Error(err))
		}
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
	exCtx := context.Background()
	if err := api.GeoResolver.Resolve(&order.Destination, exCtx); err != nil {
		api.Logger.Error("fail to resolve destination",
			zap.Error(err),
			zap.String("courier_id", *order.CourierID),
			zap.Any("dest", order.Destination))
	}
	if err := api.GeoResolver.Resolve(&order.Source, exCtx); err != nil {
		api.Logger.Error("fail to resolve source",
			zap.Error(err),
			zap.String("courier_id", *order.CourierID),
			zap.Any("dest", order.Destination))
	}
	created, err := api.OrdersDAO.Create(&order)
	if err != nil {
		api.Logger.Error("fail to create order", zap.String("courier_id", courierID), zap.Error(err))
		switch err.(type) {
		case *elastic.Error:
			err := err.(*elastic.Error)
			if err.Status == 404 {
				ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
				return
			}
		case *models.Error:
			err := err.(*models.Error)
			ctx.AbortWithStatusJSON(err.HttpStatus(), err)
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	err = api.CourierRouteDAO.CreateCourier(*order.CourierID)
	if err != nil {
		api.Logger.Error("fail to create order", zap.String("courier_id", courierID), zap.Error(err))
		//TODO: error handling
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}

	if err := api.OrdersCountTracker.Inc(ctx.Param("courier_id")); err != nil {
		api.Logger.Error("fail to increment order counter", zap.Error(err))
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
		api.Logger.Error("fail to assing order to another courier", zap.String("courier_id", courierID), zap.String("order_id", orderID), zap.Error(err))
		switch err.(type) {
		case *elastic.Error:
			err := err.(*elastic.Error)
			if err.Status == 404 {
				ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
				return
			}
		case *models.Error:
			err := err.(*models.Error)
			ctx.AbortWithStatusJSON(err.HttpStatus(), err)
			return
		}
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
		api.Logger.Error("fail to delete order", zap.String("order_id", orderID), zap.Error(err))
		switch err.(type) {
		case *elastic.Error:
			err := err.(*elastic.Error)
			if err.Status == 404 {
				ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
				return
			}
		case *models.Error:
			err := err.(*models.Error)
			ctx.AbortWithStatusJSON(err.HttpStatus(), err)
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}

	if err := api.OrdersCountTracker.Dec(ctx.Param("courier_id")); err != nil {
		api.Logger.Error("fail to decrement order counter", zap.Error(err))
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
	orders, err := api.OrdersDAO.GetOrdersForCourier(courierID,
		params.Since,
		params.Asc,
		params.ExcludeDelivered)
	if err != nil {
		api.Logger.Error("fail to get orders for courier", zap.String("courier_id", courierID), zap.Error(err))
		switch err.(type) {
		case *elastic.Error:
			err := err.(*elastic.Error)
			if err.Status == 404 {
				ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound)
				return
			}
		case *models.Error:
			err := err.(*models.Error)
			ctx.AbortWithStatusJSON(err.HttpStatus(), err)
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	ctx.JSON(http.StatusOK, orders)
}
