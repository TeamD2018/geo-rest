package controllers

import (
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"go.uber.org/zap"
)

type APIService struct {
	OrdersDAO          interfaces.IOrdersDao
	CouriersDAO        interfaces.ICouriersDAO
	CourierRouteDAO    interfaces.GeoRouteInterface
	GeoResolver        interfaces.GeoResolver
	RegionResolver     interfaces.IRegionResolver
	CourierSuggester   interfaces.CourierSuggester
	Logger             *zap.Logger
	SuggestionService  interfaces.SuggestionService
	OrdersCountTracker interfaces.OrdersCountTracker
}
