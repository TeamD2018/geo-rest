package controllers

import (
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"go.uber.org/zap"
)

type APIService struct {
	OrdersDAO        interfaces.IOrdersDao
	CouriersDAO      interfaces.ICouriersDAO
	GeoResolver      interfaces.GeoResolver
	CourierSuggester interfaces.CourierSuggester
	Logger           *zap.Logger
}
