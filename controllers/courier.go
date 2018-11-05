package controllers

import (
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"log"
	"go.uber.org/zap"
	"net/http"
	"time"
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
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	courier, err := api.CouriersDAO.GetByID(courierID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, models.ErrEntityNotFound.SetParameter(courierID))
		return
	}
	if err := api.OrdersCountTracker.Sync(models.Couriers{courier}); err != nil {
		api.Logger.Error("fail to sync courier counter", zap.Error(err))
	}
	ctx.JSON(http.StatusOK, courier)
	return
}

func (api *APIService) MiddlewareGeoSearch(ctx *gin.Context) {
	if ctx.Request.URL.Query().Get("radius") != "" {
		api.GetCouriersByCircleField(ctx)
		return
	} else {
		api.GetCouriersByBoxField(ctx)
		return
	}
}

func (api *APIService) UpdateCourier(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	if _, err := uuid.FromString(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat.SetParameter("courier_id"))
		return
	}
	courier := &models.CourierUpdate{}
	if err := ctx.BindJSON(courier); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	courier.ID = &courierID
	updated, err := api.CouriersDAO.Update(courier)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	pointWithTs := &models.PointWithTs{
		Point: updated.Location.Point,
		Ts:    uint64(time.Now().Unix()),
	}
	if orders, err := api.OrdersDAO.GetOrdersForCourier(courierID, 0, true, true); err != nil {
		api.Logger.Error("fail to retrieve counter orders", zap.Error(err), zap.String("courier_id", courierID))
	} else {
		if orders != nil && len(orders) != 0 {
			if err := api.CourierRouteDAO.AddPointToRoute(courierID, pointWithTs); err != nil {
				log.Println(err)
				//TODO: need info
			}
		}
	}
	if err := api.OrdersCountTracker.Sync(models.Couriers{updated}); err != nil {
		api.Logger.Error("fail to sync order counter", zap.Error(err), zap.String("courier_id", courierID))
	}
	ctx.JSON(http.StatusOK, updated)
	return
}

func (api *APIService) GetCouriersByCircleField(ctx *gin.Context) {
	searchParams := parameters.CircleFieldQuery{}
	if err := ctx.BindQuery(&searchParams); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	couriers, err := api.CouriersDAO.GetByCircleField(searchParams.ToCircleField(), searchParams.Size)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	if err := api.OrdersCountTracker.Sync(couriers); err != nil {
		api.Logger.Error("fail to sync couriers counters", zap.Error(err))
	}
	ctx.JSON(http.StatusOK, couriers)
	return
}

func (api *APIService) GetCouriersByBoxField(ctx *gin.Context) {
	searchParams := parameters.BoxFieldQuery{}
	if err := ctx.BindQuery(&searchParams); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	couriers, err := api.CouriersDAO.GetByBoxField(searchParams.ToBoxField(), searchParams.Size)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	if err := api.OrdersCountTracker.Sync(couriers); err != nil {
		api.Logger.Error("fail to sync couriers counters", zap.Error(err))
	}

	ctx.JSON(http.StatusOK, couriers)
	return

}

func (api *APIService) SuggestCourier(ctx *gin.Context) {
	suggestParams := parameters.Suggestion{}
	if err := ctx.BindQuery(&suggestParams); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	if suggestParams.Limit <= 0 {
		suggestParams.Limit = 200
	}
	couriers, err := api.CourierSuggester.Suggest("suggestions", &suggestParams)
	if err != nil {
		api.Logger.Error("fail to suggest couriers", zap.Error(err), zap.String("prefix", suggestParams.Prefix))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
	}
	if err := api.OrdersCountTracker.Sync(couriers); err != nil {
		api.Logger.Error("fail to sync couriers counters", zap.Error(err))
	}
	ctx.JSON(http.StatusOK, couriers)
}

func (api *APIService) GetRouteForCourier(ctx *gin.Context) {
	courierRouteParams := parameters.CourierRoute{}
	courierRouteParams.CourierID = ctx.Param("courier_id")
	if err := ctx.BindQuery(&courierRouteParams); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	if courierRouteParams.Since < 0 {
		courierRouteParams.Since = 0
	}
	if points, err := api.CourierRouteDAO.GetRoute(courierRouteParams.CourierID, courierRouteParams.Since); err != nil {
		api.Logger.Error("fail to get route", zap.Error(err),
			zap.String("courier_id", courierRouteParams.CourierID),
			zap.Int64("since", courierRouteParams.Since))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	} else {
		geoHistory := models.RouteResponse{GeoHistory: points}
		ctx.JSON(http.StatusOK, geoHistory)
	}
}

func (api *APIService) DeleteCourier(ctx *gin.Context) {
	courierID := ctx.Param("courier_id")
	if _, err := uuid.FromString(courierID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.ErrOneOfParameterHaveIncorrectFormat)
		return
	}
	if err := api.CouriersDAO.Delete(courierID); err != nil {
		api.Logger.Sugar().Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	if err := api.CourierRouteDAO.DeleteCourier(courierID); err != nil {
		api.Logger.Sugar().Error(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}
	if err := api.OrdersDAO.DeleteOrdersForCourier(courierID); err != nil {
		api.Logger.Error("fail to delete orders for courier",
			zap.Error(err),
			zap.String("courier_id", courierID))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrServerError)
		return
	}

	if err := api.OrdersCountTracker.Drop(ctx.Param("courier_id")); err != nil {
		api.Logger.Error("fail to drop orders counter", zap.Error(err), zap.String("courier_id", courierID))
	}

	ctx.Status(http.StatusNoContent)
}
