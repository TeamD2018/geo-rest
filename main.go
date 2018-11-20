package main

import (
	"flag"
	"github.com/TeamD2018/geo-rest/controllers"
	"github.com/TeamD2018/geo-rest/migrations"
	"github.com/TeamD2018/geo-rest/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"googlemaps.github.io/maps"
	"log"
	"net/http"
	"time"
)

const (
	createCouriersRouteSpaceFuncName = "create_couriers_route_space"
	createResolverCacheSpaceFuncName = "create_resolver_cache"
)

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.StringP("config", "c", "./config.toml", "path to config for geo-rest service")
	remoteConfigUrl := pflag.StringP("remote-config", "r", "", "url to config for geo-rest service")
	pflag.StringP("mode", "m", "dev", "dev/prod mode")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic("err in bind flag")
	}

	viper.SetDefault("suggestions.couriers.fuzziness", services.CouriersDefaultFuzziness)
	viper.SetDefault("suggestions.couriers.threshold", services.CouriersDefaultFuzzinessThreshold)

	if *remoteConfigUrl != "" {
		resp, err := http.Get(*remoteConfigUrl)
		if err != nil {
			panic(err)
		}

		viper.SetConfigType("toml")
		if err := viper.ReadConfig(resp.Body); err != nil {
			panic(err)
		}
	} else {
		viper.SetConfigFile(viper.GetString("config"))
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}
}

type Mode string

const (
	Production  Mode = "prod"
	Development Mode = "dev"
)

func ModeFromString(mode string) Mode {
	switch mode {
	case "prod":
		return Mode(mode)
	case "dev":
		return Mode(mode)
	default:
		panic("INVALID MODE")
	}
}

func GetLogger(mode Mode) (*zap.Logger, error) {
	switch mode {
	case Production:
		return zap.NewProduction()
	case Development:
		return zap.NewDevelopment()
	default:
		return zap.NewDevelopment()
	}
}

func main() {
	mode := ModeFromString(viper.GetString("mode"))
	logger, err := GetLogger(mode)
	if err != nil {
		log.Fatal(err)
	}
	if mode == Production {
		gin.SetMode("release")
	}
	elasticClient, err := elastic.NewClient(
		elastic.SetURL(viper.GetString("elastic.url")),
		elastic.SetSniff(viper.GetBool("elastic.sniff")),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewConstantBackoff(time.Second*5))))
	if err != nil {
		log.Fatal(err)
	}
	gmaps, err := maps.NewClient(maps.WithAPIKey(viper.GetString("google-maps.apikey")))
	if err != nil {
		log.Fatal(err)
	}
	tntClient, err := tarantool.Connect(viper.GetString("tarantool.url"), tarantool.Opts{
		User: viper.GetString("tarantool.user"),
		Pass: viper.GetString("tarantool.pass"),
	})
	if err != nil {
		log.Fatal(err)
	}
	err = migrations.Driver{Client: tntClient, Logger: logger}.Run()
	if err != nil {
		logger.Fatal("fail to perform migrations", zap.Error(err))
	}
	_, err = tntClient.Call17(createCouriersRouteSpaceFuncName, []interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	_, err = tntClient.Call17(createResolverCacheSpaceFuncName, []interface{}{})
	if err != nil {
		log.Fatal(err)
	}

	ordersCountTracker := services.NewTarantoolOrdersCountTracker(tntClient, logger)

	couriersDao := services.NewCouriersElasticDAO(elasticClient, logger, "", services.DefaultCouriersReturnSize)
	ordersDao := services.NewOrdersElasticDAO(elasticClient, logger, couriersDao, "")

	tntResolver := services.NewTntResolver(tntClient, logger)
	gmapsResolver := services.NewGMapsResolver(gmaps, logger)

	couriersSuggester := services.NewCouriersSuggesterElastic(elasticClient, couriersDao, logger).
		SetFuzziness(viper.GetInt("suggestions.couriers.fuzziness")).
		SetFuzzinessThreshold(viper.GetInt("suggestions.couriers.threshold"))

	couriersSuggestEngine := services.PrefixSuggestEngine{
		Fuzziness: "AUTO",
		Limit:     15,
		Field:     "suggestions",
		Index:     couriersDao.GetIndex(),
	}
	ordersPrefixSuggestEngine := services.PrefixSuggestEngine{
		Fuzziness: "0",
		Limit:     15,
		Field:     "order_suggestions",
		Index:     ordersDao.GetIndex(),
	}
	ordersSuggestDestinationEngine := services.OrdersSuggestEngine{
		Fuzziness:          "1",
		FuzzinessThreshold: 5,
		Limit:              15,
		Field:              "destination.address",
		Index:              ordersDao.GetIndex(),
	}

	suggestersExecutor := services.NewSuggestEngineExecutor(elasticClient, logger)
	suggestersExecutor.AddEngine("orders-engine", &ordersSuggestDestinationEngine)
	suggestersExecutor.AddEngine("couriers-engine", &couriersSuggestEngine)
	suggestersExecutor.AddEngine("orders-prefix-engine", &ordersPrefixSuggestEngine)

	if err := couriersDao.EnsureMapping(); err != nil {
		logger.Fatal("Fail to ensure couriers mapping: ", zap.Error(err))
	}
	if err := ordersDao.EnsureMapping(); err != nil {
		logger.Fatal("Fail to ensure orders mapping: ", zap.Error(err))
	}

	tntRouteDao := services.NewTarantoolRouteDAO(tntClient, logger)

	api := controllers.APIService{
		CouriersDAO:        couriersDao,
		OrdersDAO:          ordersDao,
		CourierRouteDAO:    tntRouteDao,
		GeoResolver:        services.NewCachedResolver(tntResolver, gmapsResolver),
		CourierSuggester:   couriersSuggester,
		Logger:             logger,
		SuggesterExecutor:  suggestersExecutor,
		OrdersCountTracker: ordersCountTracker,
	}
	router := gin.New()

	router.Use(func(ctx *gin.Context) {
		ctx.Set(controllers.LoggerKey, logger)
	}, controllers.LogBody)

	pprof.Register(router)

	config := cors.DefaultConfig()
	config.AllowOrigins = viper.GetStringSlice("cors.origins")
	router.Use(cors.New(config))

	controllers.SetupRouters(router, &api)

	if err := router.Run(viper.GetString("server.url")); err != nil {
		log.Fatal(err)
	}
}
