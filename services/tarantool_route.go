package services

import (
	"github.com/olivere/elastic"
	"github.com/tarantool/go-tarantool"
)

const (
	addCourierWithOrderFuncName = "add_courier_with_order"
	addPointToRouteFuncName     = "add_point_to_route"
	getRouteFuncName            = "get_route"
)

const (
	deleteCourierFuncName = "delete_courier"
)

type TarantoolRouteDAO struct {
	client *tarantool.Connection
}

func (tnt *TarantoolRouteDAO) DeleteCourier(courierID string) error {
	_, err := tnt.client.Call17(deleteCourierFuncName, []interface{}{courierID})
	if err != nil {
		return err
	}
	return nil
}

func NewTarantoolRouteDAO(client *tarantool.Connection) *TarantoolRouteDAO {
	return &TarantoolRouteDAO{client: client}
}

func (tnt *TarantoolRouteDAO) CreateCourierWithOrder(courierID, orderID string) error {
	_, err := tnt.client.Call17(addCourierWithOrderFuncName, []interface{}{courierID, orderID})
	if err != nil {
		return err
	}
	return nil
}

func (tnt *TarantoolRouteDAO) AddPointToRoute(courierID string, point *elastic.GeoPoint) error {
	p := map[string]float64{
		"lat": point.Lat,
		"lon": point.Lon,
	}
	_, err := tnt.client.Call17(addPointToRouteFuncName, []interface{}{courierID, p})
	if err != nil {
		return err
	}
	return nil
}

func (tnt *TarantoolRouteDAO) GetRoute(courierID, orderID string) ([]*elastic.GeoPoint, error) {
	resp, err := tnt.client.Call17(getRouteFuncName, []interface{}{courierID, orderID})
	points := make([]*elastic.GeoPoint, len(resp.Data))
	if err != nil {
		return nil, err
	}
	for i, p := range resp.Data {
		points[i] = elastic.GeoPointFromLatLon(p.(map[interface{}]interface{})["lat"].(float64), p.(map[interface{}]interface{})["lon"].(float64))
	}
	return points, nil
}
