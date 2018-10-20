package services

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

const (
	addCourierWithOrderFuncName = "add_courier"
	addPointToRouteFuncName     = "add_point_to_route"
	getRouteFuncName            = "get_route"
)

const (
	deleteCourierFuncName = "delete_courier"
)

type TarantoolRouteDAO struct {
	l      *zap.Logger
	client *tarantool.Connection
}

func (tnt *TarantoolRouteDAO) DeleteCourier(courierID string) error {
	_, err := tnt.client.Call17(deleteCourierFuncName, []interface{}{courierID})
	if err != nil {
		return err
	}
	return nil
}

func NewTarantoolRouteDAO(client *tarantool.Connection, logger *zap.Logger) *TarantoolRouteDAO {
	return &TarantoolRouteDAO{
		client: client,
		l:      logger,
	}
}

func (tnt *TarantoolRouteDAO) CreateCourier(courierID string) error {
	_, err := tnt.client.Call17(addCourierWithOrderFuncName, []interface{}{courierID})
	if err != nil {
		tnt.l.Sugar().Error(err)
		return err
	}
	return nil
}

func (tnt *TarantoolRouteDAO) AddPointToRoute(courierID string, point *models.PointWithTs) error {
	p := map[string]interface{}{
		"lat": point.Point.Lat,
		"lon": point.Point.Lon,
		"ts":  point.Ts,
	}
	_, err := tnt.client.Call17(addPointToRouteFuncName, []interface{}{courierID, p})
	if err != nil {
		tnt.l.Sugar().Error(err)
		return err
	}
	return nil
}

func (tnt *TarantoolRouteDAO) GetRoute(courierID string) ([]*models.PointWithTs, error) {
	resp, err := tnt.client.Call17(getRouteFuncName, []interface{}{courierID})
	points := make([]*models.PointWithTs, len(resp.Data))
	if err != nil {
		tnt.l.Sugar().Errorw("msg", "resp", resp, "error", err)
		return nil, err
	}
	for i, p := range resp.Data {
		points[i] = &models.PointWithTs{
			Point: elastic.GeoPointFromLatLon(p.(map[interface{}]interface{})["lat"].(float64), p.(map[interface{}]interface{})["lon"].(float64)),
			Ts:    p.(map[interface{}]interface{})["ts"].(uint64),
		}
	}
	return points, nil
}
