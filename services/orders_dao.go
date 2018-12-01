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
	"strconv"
	"time"
)

const OrdersIndex = "order"

type OrdersElasticDAO struct {
	Elastic     *elastic.Client
	couriersDAO interfaces.ICouriersDAO
	index       string
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
		index:       index,
		Logger:      logger,
		couriersDAO: couriersDAO,
	}
}

func (od *OrdersElasticDAO) Get(orderID string) (*models.Order, error) {
	db := od.Elastic
	orderRaw, err := db.Get().
		Index(od.index).
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
	var order orderWrapper
	order.Source = orderCreate.Source
	order.Destination = orderCreate.Destination
	order.CourierID = *orderCreate.CourierID
	order.OrderNumber = orderCreate.OrderNumber
	order.CreatedAt = time.Now().Unix()
	order.Suggestions = elastic.NewSuggestField(strconv.Itoa(order.OrderNumber))
	id := uuid.NewV4().String()
	ret, err := db.Index().
		Index(od.index).
		Type("_doc").
		Id(id).
		BodyJson(order).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	order.ID = ret.Id
	return &order.Order, nil
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
		Index(od.index).
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
		Index(od.index).
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

func (od *OrdersElasticDAO) DeleteOrdersForCourier(courierID string) error {
	db := od.Elastic
	courierIDQuery := elastic.NewTermQuery("courier_id", courierID)
	res, err := db.DeleteByQuery(od.GetIndex()).Query(courierIDQuery).Do(context.Background())
	if err != nil {
		if res != nil {
			od.Logger.Warn("fail to delete orders", zap.Any("response", res))
		}
		if elastic.IsNotFound(err) {
			return models.ErrEntityNotFound.SetParameter(courierID)
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
	res, err := db.Search(od.index).Type("_doc").Query(ordersQuery).Do(context.Background())
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

	if !exists {
		_, err := od.Elastic.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (od *OrdersElasticDAO) GetIndex() string {
	return od.index
}

func (od *OrdersElasticDAO) GetMapping() (indexName string, mapping string) {
	return od.index, `{
 "settings": {
    "analysis": {
      "analyzer": {
        "autocomplete": {
          "tokenizer": "autocomplete",
          "filter": [
            "lowercase",
            "split_words",
            "remove_frequent_words",
            "synonym",
            "min_letter_trigram_len",
            "discard_empty_strings",
            "unique"
          ]
        },
        "autocomplete_search": {
          "tokenizer": "whitespace",
          "filter": [
            "lowercase",
            "split_words",
            "min_letter_trigram_len",
            "discard_empty_strings",
            "remove_frequent_words"
          ]
        }
      },
      "filter": {
        "split_words": {
          "type": "word_delimiter",
          "preserve_original": false
        },
        "synonym": {
          "type": "synonym",
          "synonyms": [
            "стр => строение",
            "корп => корпус",
            "кв => квартира",
            "пр => проспект",
            "пер => переулок",
            "пл => площадь"
          ]
        },
        "remove_frequent_words": {
          "type": "stop",
          "stopwords": [
            "улица",
            "yл.",
            "ул",
            "дом"
          ]
        },
        "min_letter_trigram_len": {
          "type": "pattern_replace",
          "pattern": "^\\D{1,2}$"
        },
        "discard_empty_strings": {
          "type": "length",
          "min": 1
        },
        "min_token_length": {
          "min": 3,
          "type": "length"
        }
      },
      "tokenizer": {
        "autocomplete": {
          "type": "edge_ngram",
          "min_gram": 1,
          "max_gram": 15,
          "token_chars": [
            "letter",
            "digit"
          ]
        }
      }
    }
  },
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
		"order_number": {
          "type": "integer"
        },
		"order_suggestions": {
		  "type": "completion",
		  "analyzer": "whitespace"
		},
        "destination": {
          "properties": {
            "point": {
              "type": "geo_point"
            },
            "address": {
              "type": "text",
              "analyzer": "autocomplete",
              "search_analyzer": "autocomplete_search"
            }
          }
        },
        "source": {
          "properties": {
            "point": {
              "type": "geo_point"
            },
            "address": {
              "type": "text",
              "analyzer": "autocomplete",
              "search_analyzer": "autocomplete_search"
            }
          }
        }
      }
    }
  }
}`
}

type orderWrapper struct {
	Suggestions *elastic.SuggestField `json:"order_suggestions,omitempty"`
	models.Order
}
