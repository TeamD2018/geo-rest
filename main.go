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
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

const (
	createCouriersRouteSpaceFuncName = "create_couriers_route_space"
	createResolverCacheSpaceFuncName = "create_resolver_cache"
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
	tntClient, err := tarantool.Connect(viper.GetString("tarantool.url"), tarantool.Opts{})
	if err != nil {
		log.Fatal(err)
	}
	code := getLuaCode()
	_, err = tntClient.Eval(code, []interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	_, err = tntClient.Call17(createCouriersRouteSpaceFuncName, []interface{}{})
	_, err = tntClient.Call17(createResolverCacheSpaceFuncName, []interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	couriersDao := services.NewCouriersElasticDAO(elasticClient, logger, "", services.DefaultCouriersReturnSize)
	ordersDao := services.NewOrdersElasticDAO(elasticClient, logger, couriersDao, "")

	tntResolver := services.NewTntResolver(tntClient, logger)
	gmapsResolver := services.NewGMapsResolver(gmaps, logger)
	if err := couriersDao.EnsureMapping();
		err != nil {
		logger.Fatal("Fail to ensure couriers mapping: ", zap.Error(err))
	}
	if err := ordersDao.EnsureMapping(); err != nil {
		logger.Fatal("Fail to ensure orders mapping: ", zap.Error(err))
	}

	tntRouteDao := services.NewTarantoolRouteDAO(tntClient)

	api := controllers.APIService{
		CouriersDAO:     couriersDao,
		OrdersDAO:       ordersDao,
		CourierRouteDAO: tntRouteDao,
		GeoResolver:     services.NewCachedResolver(tntResolver, gmapsResolver),
		Logger:          logger,
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

func getLuaCode() string {
	dirname := "./tnt_stored_procedures/"
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	buf := strings.Builder{}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".lua" {
			code, err := ioutil.ReadFile(dirname + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			buf.Write(code)
		}
	}
	return buf.String()
}
