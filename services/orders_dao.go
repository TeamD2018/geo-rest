package services

import (
	"context"
	"encoding/json"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

const OrdersIndex = "order"

type OrdersElasticDAO struct {
	Elastic *elastic.Client
	Index   string
	Logger  *zap.Logger
}

func NewOrdersElasticDAO(client *elastic.Client, logger *zap.Logger, index string) *OrdersElasticDAO {
	if index == "" {
		index = OrdersIndex
	}
	return &OrdersElasticDAO{
		Elastic: client,
		Index:   index,
		Logger:  logger,
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
	if err := json.Unmarshal(*orderRaw.Source, &order); err != nil {
		return nil, err
	}
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
	id := *update.ID
	update.ID = nil
	orderRaw, err := db.Update().
		Index(od.Index).
		Id(id).
		Doc(*update).
		FetchSource(true).
		Do(context.Background())
	if err != nil {
		od.Logger.Sugar().Errorw("order update failed", *update)
		return nil, err
	}
	var order models.Order
	if err := json.Unmarshal(*orderRaw.GetResult.Source, &order); err != nil {
		return nil, err
	}
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

func (od *OrdersElasticDAO) GetMapping() (indexName string, mapping string) {
	return "order", `{
		"mappings": {
			"_doc": {
				"properties": {
					"courier_id": {
						"type": "keyword"
					},
					"created_at": {
						"type": "long"
					},
					"delivered_at": {
						"type": "long"
					},
					"destination": {
						"properties": {
							"geo_point": {
								"type": "geo_point"
							},
							"address": {
								"type": "completion"
							}
						}
					},
					"source": {
						"properties": {
							"geo_point": {
								"type": "geo_point"
							},
							"address": {
								"type": "completion"
							}
						}
					}
				}	
			}
		}		
	}`
}
