package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

const CourierIndex = "couriers"

type CouriersElasticDAO struct {
	client *elastic.Client
	index  string
	l      *zap.Logger
}

func NewCouriersElasticDAO(client *elastic.Client, logger *zap.Logger, index string) *CouriersElasticDAO {
	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}
	if index == "" {
		index = CourierIndex
	}
	return &CouriersElasticDAO{client: client, index: index, l: logger}
}

func (c *CouriersElasticDAO) GetByID(courierID string) (*models.Courier, error) {
	res, err := c.client.Get().Index(c.index).Type("_doc").Id(courierID).Do(context.Background())
	if err != nil {
		c.l.Sugar().Errorw("", zap.Error(err))
		if err.(*elastic.Error).Status == 404 {
			return nil, models.ErrCourierNotFound.SetParameter(courierID)
		}
		return nil, err
	}
	c.l.Sugar().Debug(res)
	result := &models.Courier{}
	if err := json.Unmarshal(*res.Source, &result); err != nil {
		c.l.Sugar().Errorw("", zap.Error(err))
		return nil, models.ErrUnmarshalJSON.SetParameter(err)
	}
	result.ID = res.Id
	return result, nil
}

func (*CouriersElasticDAO) GetByName(name string) (models.Couriers, error) {
	return nil, nil
}

func (*CouriersElasticDAO) GetBySquareField(field *models.SquareField) (models.Couriers, error) {
	return nil, nil
}

func (c *CouriersElasticDAO) GetByCircleField(field *models.CircleField) (models.Couriers, error) {
	boolQuery := elastic.NewBoolQuery()
	geodistanceQuery := elastic.NewGeoDistanceQuery("location.point").
		GeoPoint(field.Center).
		Distance(fmt.Sprintf("%dm", field.Radius))
	match := elastic.NewMatchAllQuery()

	query := boolQuery.Must(match).
		Filter(geodistanceQuery)

	end := c.client.Search(c.index).
		Type("_doc").
		Query(query)

	result := models.Couriers{}
	res, err := end.Do(context.Background())
	if err != nil {
		return nil, err
	}

	for _, item := range res.Hits.Hits {
		var courier models.Courier
		if err := json.Unmarshal(*item.Source, &courier);
		err != nil {
			return nil, err
		}
		courier.ID = item.Id
		result = append(result, &courier)
	}
	return result, nil
}


func (c *CouriersElasticDAO) Create(courier *models.CourierCreate) (*models.Courier, error) {
	m := &models.Courier{
		Name:  courier.Name,
		Phone: courier.Phone,
	}
	id := uuid.NewV4().String()
	res, err := c.client.Index().
		Index(c.index).
		Type("_doc").
		Id(id).
		BodyJson(m).
		Do(context.Background())
	if err != nil {
		c.l.Sugar().Errorw("", zap.Error(err))
		return nil, err
	}
	m.ID = id
	c.l.Sugar().Debugw("", zap.Any("res", res))
	return m, nil
}

func (c *CouriersElasticDAO) Update(courier *models.CourierUpdate) (*models.Courier, error) {
	id := *courier.ID
	courier.ID = nil
	if courier.Location != nil {
		now := time.Now().Unix()
		courier.LastSeen = &now
	}


	res, err := c.client.Update().
		Index(c.index).
		Type("_doc").
		Id(id).
		Doc(courier).
		FetchSource(true).
		Do(context.Background())
	if err != nil {
		c.l.Sugar().Error(zap.Error(err))
		return nil, err
	}
	result := &models.Courier{}
	if err := json.Unmarshal(*res.GetResult.Source, &result); err != nil {
		c.l.Sugar().Errorw("", zap.Error(err))
		return nil, models.ErrUnmarshalJSON.SetParameter(err)
	}
	c.l.Sugar().Debugw("", zap.Any("res", res))
	result.ID = res.Id
	return result, nil
}

func (c *CouriersElasticDAO) Delete(courierID string) error {
	res, err := c.client.Delete().Index(c.index).Type("_doc").Id(courierID).Do(context.Background())
	if err != nil {
		if err.(*elastic.Error).Status == 404 {
			return models.ErrCourierNotFound.SetParameter(courierID)
		}
		c.l.Sugar().Errorw("", zap.Error(err))
		return err
	}
	c.l.Sugar().Debugw("", zap.Any("res", res))
	return nil
}

func (c *CouriersElasticDAO) EnsureMapping() error {
	indexName, mapping := c.GetMapping()

	ctx := context.Background()
	exists, err := c.client.IndexExists(indexName).Do(ctx)
	if err != nil {
		c.l.Sugar().Errorw("", zap.Error(err))
		return err
	}

	if exists == false {
		_, err := c.client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			c.l.Sugar().Errorw("", zap.Error(err))
			return err
		}
	}

	return nil
}

func (c *CouriersElasticDAO) GetMapping() (indexName string, mapping string) {
	return c.index, `{
		"mappings": {
			"_doc": {
				"properties": {
					"name": {
						"type": "keyword"
					},
					"location": {
						"properties": {
							"point": {
								"type": "geo_point"
							},
							"address": {
								"type": "completion"
							}
						}
					},
					"phone": {
						"type": "keyword",
						"index": false
					},
					"last_seen": {
						"type": "long",
						"index": false
					}
				}
			}
		}		
	}`
}
