package controllers

import (
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"go.uber.org/zap"
)

type APIService struct {
	OrdersDAO       interfaces.IOrdersDao
	CouriersDAO     interfaces.ICouriersDAO
	CourierRouteDAO interfaces.GeoRouteInterface
	GeoResolver     interfaces.GeoResolver
	Logger          *zap.Logger
}
