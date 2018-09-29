package controllers

import (
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"go.uber.org/zap"
)

type ApiService struct {
	OrdersDAO interfaces.IOrdersDao
	Logger    *zap.Logger
}
