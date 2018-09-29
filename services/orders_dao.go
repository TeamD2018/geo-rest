package services

import (
	"context"
	"encoding/json"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/satori/go.uuid"
	"time"
)

const OrdersIndex = "order"

type OrdersElasticDAO struct {
	Elastic *elastic.Client
	Index   string
}

func NewOrdersElasticDAO(client *elastic.Client, index string) *OrdersElasticDAO {
	if index == "" {
		index = OrdersIndex
	}
	return &OrdersElasticDAO{
		Elastic: client,
		Index:   index,
	}
}

func (od *OrdersElasticDAO) Get(orderID string) (*models.Order, error) {
	db := od.Elastic
	orderRaw, err := db.Get().
		Index(od.Index).
		Id(orderID).
		Do(context.Background())

	var order models.Order
	if err != nil {
		return nil, err
	}
	json.Unmarshal(*orderRaw.Source, &order)
	return &order, nil
}

func (od *OrdersElasticDAO) Create(orderCreate *models.OrderCreate) (*models.Order, error) {
	db := od.Elastic
	var order models.Order
	order.Source = orderCreate.Source
	order.Destination = orderCreate.Destination
	order.CourierID = *orderCreate.CourierID
	order.CreatedAt = time.Now().Unix()
	order.ID = uuid.NewV4().String()
	_, err := db.Index().
		Index(od.Index).
		BodyJson(order).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (od *OrdersElasticDAO) Update(update *models.OrderUpdate) (*models.Order, error) {
	db := od.Elastic
	orderRaw, err := db.Update().
		Index(od.Index).
		Id(*update.ID).
		Doc(*update).
		FetchSource(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	var order models.Order
	json.Unmarshal(*orderRaw.GetResult.Source, &order)
	return &order, nil
}

func (od *OrdersElasticDAO) Delete(orderID string) error {
	db := od.Elastic
	_, err := db.Delete().
		Index(od.Index).
		Id(orderID).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
