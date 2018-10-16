package main

import (
	"flag"
	"github.com/TeamD2018/geo-rest/controllers"
	"github.com/TeamD2018/geo-rest/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"log"
)

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.StringP("config", "c", "./config.toml", "path to config for geo-rest service")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic("err in bind flag")
	}
	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func main() {
	elasticClient, err := elastic.NewClient(
		elastic.SetURL(viper.GetString("elastic.url")),
		elastic.SetSniff(viper.GetBool("elastic.sniff")))
	if err != nil {
		log.Fatal(err)
	}
	gmaps, err := maps.NewClient(maps.WithAPIKey(viper.GetString("google-maps.apikey")))
	if err != nil {
		log.Fatal(err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	couriersDao := services.NewCouriersElasticDAO(elasticClient, logger, "", services.DefaultCouriersReturnSize)
	ordersDao := services.NewOrdersElasticDAO(elasticClient, logger, couriersDao, "")
	gmapsResolver := services.NewGMapsResolver(gmaps, logger)
	couriersSuggester := services.NewCouriersSuggesterElastic(elasticClient, couriersDao, logger)
	if err := couriersDao.EnsureMapping();
		err != nil {
		logger.Fatal("Fail to ensure couriers mapping: ", zap.Error(err))
	}
	if err := ordersDao.EnsureMapping();
		err != nil {
		logger.Fatal("Fail to ensure orders mapping: ", zap.Error(err))
	}

	api := controllers.APIService{
		CouriersDAO:      couriersDao,
		OrdersDAO:        ordersDao,
		GeoResolver:      gmapsResolver,
		CourierSuggester: couriersSuggester,
		Logger:           logger,
	}
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = viper.GetStringSlice("cors.origins")
	router.Use(cors.New(config))

	controllers.SetupRouters(router, &api)

	if err := router.Run(viper.GetString("server.url")); err != nil {
		log.Fatal(err)
	}
}
