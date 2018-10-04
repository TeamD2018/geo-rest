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

	controllers.SetupRouters(router, &api)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
