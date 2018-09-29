package main

import (
	"github.com/TeamD2018/geo-rest/controllers"
	"github.com/TeamD2018/geo-rest/services"
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
	api := controllers.APIService{
		OrdersDAO: services.NewOrdersElasticDAO(elasticClient, logger, ""),
		Logger:    logger,
	}
	router := gin.Default()
	router.POST("/couriers/:courier_id/orders", api.CreateOrder)
	router.GET("/couriers/:courier_id/orders/:order_id", api.GetOrder)
	router.PUT("/couriers/:courier_id/orders/:order_id", api.UpdateOrder)
	router.PATCH("/couriers/:courier_id/orders/:order_id", api.AssignNewCourier)
	router.DELETE("/couriers/:courier_id/orders/:order_id", api.DeleteOrder)
	router.Run(":8080")
}
