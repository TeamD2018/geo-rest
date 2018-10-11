package services

import (
	"context"
	"encoding/json"
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/olivere/elastic"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

const OrdersIndex = "order"

type OrdersElasticDAO struct {
	Elastic     *elastic.Client
	couriersDAO interfaces.ICouriersDAO
	Index       string
	Logger      *zap.Logger
}

func NewOrdersElasticDAO(client *elastic.Client,
	logger *zap.Logger,
	couriersDAO interfaces.ICouriersDAO,
	index string) *OrdersElasticDAO {
	if index == "" {
		index = OrdersIndex
	}
	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}
	return &OrdersElasticDAO{
		Elastic:     client,
		Index:       index,
		Logger:      logger,
		couriersDAO: couriersDAO,
	}
}

func (od *OrdersElasticDAO) Get(orderID string) (*models.Order, error) {
	db := od.Elastic
	orderRaw, err := db.Get().
		Index(od.Index).
		Type("_doc").
		Id(orderID).
		Do(context.Background())

	var order models.Order
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, models.ErrEntityNotFound.SetParameter(orderID)
		}
		return nil, err
	}
	if err := json.Unmarshal(*orderRaw.Source, &order); err != nil {
		return nil, models.ErrUnmarshalJSON
	}
	order.ID = orderRaw.Id
	return &order, nil
}

func (od *OrdersElasticDAO) Create(orderCreate *models.OrderCreate) (*models.Order, error) {
	db := od.Elastic
	exists, err := od.couriersDAO.Exists(*orderCreate.CourierID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, models.ErrEntityNotFound.SetParameter(*orderCreate.CourierID)
	}
	var order models.Order
	order.Source = orderCreate.Source
	order.Destination = orderCreate.Destination
	order.CourierID = *orderCreate.CourierID
	order.CreatedAt = time.Now().Unix()
	id := uuid.NewV4().String()
	ret, err := db.Index().
		Index(od.Index).
		Type("_doc").
		Id(id).
		BodyJson(order).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	order.ID = ret.Id
	return &order, nil
}

func (od *OrdersElasticDAO) Update(update *models.OrderUpdate) (*models.Order, error) {
	db := od.Elastic
	id := *update.ID
	update.ID = nil
	if update.CourierID != nil {
		exists, err := od.couriersDAO.Exists(*update.CourierID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, models.ErrEntityNotFound.SetParameter(*update.CourierID)
		}
	}
	orderRaw, err := db.Update().
		Index(od.Index).
		Type("_doc").
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
	order.ID = orderRaw.Id
	return &order, nil
}

func (od *OrdersElasticDAO) Delete(orderID string) error {
	db := od.Elastic
	_, err := db.Delete().
		Index(od.Index).
		Type("_doc").
		Id(orderID).
		Do(context.Background())
	if err != nil {
		if elastic.IsNotFound(err) {
			return models.ErrEntityNotFound.SetParameter(orderID)
		}
		return err
	}
	return nil
}

func (od *OrdersElasticDAO) GetOrdersForCourier(
	courierID string,
	since int64,
	isLowerThreshold parameters.DirectionFlag,
	excludeDelivered parameters.DeliveredFlag) (models.Orders, error) {

	db := od.Elastic
	courierIDQuery := elastic.NewTermQuery("courier_id", courierID)
	var sinceRangeQuery elastic.Query
	if isLowerThreshold {
		sinceRangeQuery = elastic.NewRangeQuery("created_at").Gte(since).IncludeUpper(false)
	} else {
		sinceRangeQuery = elastic.NewRangeQuery("created_at").Lte(since).IncludeLower(false)
	}
	ordersQuery := elastic.NewBoolQuery().Filter(courierIDQuery, sinceRangeQuery)
	if excludeDelivered {
		ordersQuery = ordersQuery.MustNot(elastic.NewExistsQuery("delivered_at"))
	}
	res, err := db.Search(od.Index).Type("_doc").Query(ordersQuery).Do(context.Background())
	if err != nil {
		return nil, err
	}
	orders := make(models.Orders, 0, res.Hits.TotalHits)
	for _, hit := range res.Hits.Hits {
		var order models.Order
		if err := json.Unmarshal(*hit.Source, &order); err != nil {
			return nil, err
		}
		order.ID = hit.Id
		orders = append(orders, &order)
	}
	return orders, nil

}

func (od *OrdersElasticDAO) EnsureMapping() error {
	indexName, mapping := od.GetMapping()

	ctx := context.Background()
	exists, err := od.Elastic.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}

	if exists == false {
		_, err := od.Elastic.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (od *OrdersElasticDAO) GetMapping() (indexName string, mapping string) {
	return od.Index, `{
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
							"point": {
								"type": "geo_point"
							},
							"address": {
								"type": "completion"
							}
						}
					},
					"source": {
						"properties": {
							"point": {
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
