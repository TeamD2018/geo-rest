package services

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

type TarantoolOrdersCountTracker struct {
	logger *zap.Logger
	db     *tarantool.Connection
}

const (
	IncCourierOrdersCount  = "inc_courier_orders_counter"
	DecCourierOrdersCount  = "dec_courier_orders_counter"
	GetCounters            = "get_counters"
	DropCourierOrdersCount = "drop_courier_orders_counter"
)

func NewTarantoolOrdersCountTracker(con *tarantool.Connection, logger *zap.Logger) *TarantoolOrdersCountTracker {
	return &TarantoolOrdersCountTracker{
		db:     con,
		logger: logger,
	}
}

func (oct *TarantoolOrdersCountTracker) Inc(courierId string) error {
	db := oct.db
	_, err := db.Call17(IncCourierOrdersCount, []interface{}{courierId})
	return err
}

func (oct *TarantoolOrdersCountTracker) Dec(courierId string) error {
	db := oct.db
	_, err := db.Call17(DecCourierOrdersCount, []interface{}{courierId})
	return err
}

func (oct *TarantoolOrdersCountTracker) Sync(couriers models.Couriers) error {
	ids := make([]string, 0, len(couriers))
	for _, courier := range couriers {
		ids = append(ids, courier.ID)
	}
	db := oct.db
	res, err := db.Call17(GetCounters, []interface{}{ids})
	if err != nil {
		return err
	}
	counters := make(map[string]int)
	oct.logger.Sugar().Debugw("got", "data", res.Data)
	for _, rawCounter := range res.Data[0].([]interface{}) {
		counter := rawCounter.([]interface{})
		courierID := counter[0].(string)
		switch total := counter[1].(type) {
		case uint64:
			counters[courierID] = int(total)
		case int64:
			counters[courierID] = int(total)
		case int:
			counters[courierID] = total
		default:
			counters[courierID] = total.(int)
		}
	}
	for _, courier := range couriers {
		courier.OrdersCount = counters[courier.ID]
	}
	return nil
}

func (oct *TarantoolOrdersCountTracker) Drop(courierId string) error {
	db := oct.db
	_, err := db.Call17(DropCourierOrdersCount, []interface{}{courierId})
	return err
}
