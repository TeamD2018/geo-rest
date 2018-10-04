package main

import (
	"github.com/TeamD2018/geo-rest/controllers"
	"github.com/TeamD2018/geo-rest/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"log"
)

func main() {
	elasticClient, err := elastic.NewClient(elastic.SetURL("http://elastic:9200"))
	if err != nil {
		log.Fatal(err)
	}
	gmaps, err := maps.NewClient(maps.WithAPIKey("AIzaSyCsvYa45nNh7NNLE_PUix8SOI73_HlcTX8"))
	if err != nil {
		log.Fatal(err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	couriersDao := services.NewCouriersElasticDAO(elasticClient, logger, "")
	ordersDao := services.NewOrdersElasticDAO(elasticClient, logger, couriersDao, "")
	gmapsResolver := services.NewGMapsResolver(gmaps, logger)
	if err := couriersDao.EnsureMapping();
		err != nil {
		logger.Fatal("Fail to ensure couriers mapping: ", zap.Error(err))
	}
	if err := ordersDao.EnsureMapping();
		err != nil {
		logger.Fatal("Fail to ensure orders mapping: ", zap.Error(err))
	}

	api := controllers.APIService{
		CouriersDAO: couriersDao,
		OrdersDAO:   ordersDao,
		GeoResolver: gmapsResolver,
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
