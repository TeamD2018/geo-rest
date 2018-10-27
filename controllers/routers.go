package controllers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouters(router *gin.Engine, api *APIService) {
	g := router.Group(`/couriers`)
	//orders endpoints
	g.POST("/:courier_id/orders", api.CreateOrder)
	g.GET("/:courier_id/orders/:order_id", api.GetOrder)
	g.PUT("/:courier_id/orders/:order_id", api.UpdateOrder)
	g.PATCH("/:courier_id/orders/:order_id", api.AssignNewCourier)
	g.DELETE("/:courier_id/orders/:order_id", api.DeleteOrder)
	g.GET("/:courier_id/orders", api.GetOrdersForCourier)

	//couriers endpoints
	g.POST("", api.CreateCourier)
	g.GET("", api.MiddlewareGeoSearch)
	g.GET("/:courier_id", api.GetCourierByID)
	g.PUT("/:courier_id", api.UpdateCourier)
	g.DELETE("/:courier_id", api.DeleteCourier)
	g.GET("/:courier_id/geo_history", api.GetRouteForCourier)

	router.GET("/suggestions/couriers", api.SuggestCourier)
	router.GET("/suggestions", api.Suggest)
}
