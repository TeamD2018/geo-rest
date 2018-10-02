package main

import (
	"github.com/TeamD2018/geo-rest/controllers"
	"github.com/TeamD2018/geo-rest/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"log"
)

func main() {
	elasticClient, err := elastic.NewClient(elastic.SetURL("http://elastic:9200"))
	if err != nil {
		log.Fatal(err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	couriersDao := services.NewCouriersElasticDAO(elasticClient, logger, "")
	ordersDao := services.NewOrdersElasticDAO(elasticClient, logger, couriersDao, "")
	api := controllers.APIService{
		CouriersDAO: couriersDao,
		OrdersDAO:   ordersDao,
		Logger:      logger,
	}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://dc.utkin.xyz:8080", "http://35.204.198.186:8080"}
	router.Use(cors.New(config))

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
	g.GET("/:courier_id", api.GetCourierByID)
	g.PUT("/:courier_id", api.UpdateCourier)
	g.DELETE("/:courier_id", api.DeleteCourier)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
