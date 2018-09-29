package services

import (
	"context"
	"encoding/json"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
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

func NewCouriersDAO(client *elastic.Client, index string) *CouriersElasticDAO {
	return &CouriersElasticDAO{client: client, index: index}
}

func (c *CouriersElasticDAO) GetByID(courierID string) (*models.Courier, error) {
	res, err := c.client.Get().Index(c.index).Type("_doc").Do(context.Background())
	if err != nil {
		c.l.Sugar().Error(zap.Error(err))
		return nil, err
	}
	c.l.Sugar().Debug(res)
	result := &models.Courier{}
	if err := json.Unmarshal(*res.Source, &result); err != nil {
		c.l.Sugar().Error(zap.Error(err))
		return nil, errors.Errorf("error due to unmarshaling json: %s", err)
	}
	return result, nil
}

func (*CouriersElasticDAO) GetByName(name string) (models.Couriers, error) {
	return nil, nil
}

func (*CouriersElasticDAO) GetBySquareField(field *models.SquareField) (*models.Courier, error) {
	return nil, nil
}

func (*CouriersElasticDAO) GetByCircleField(field *models.CircleField) (*models.Courier, error) {
	return nil, nil
}

func (c *CouriersElasticDAO) Create(courier *models.CourierCreate) (*models.Courier, error) {
	m := &models.Courier{
		ID:    uuid.NewV4().String(),
		Name:  courier.Name,
		Phone: courier.Phone,
	}
	res, err := c.client.Index().
		Index(c.index).
		Type("_doc").
		Id(m.ID).
		BodyJson(m).
		Do(context.Background())
	if err != nil {
		c.l.Sugar().Error(zap.Error(err))
		return nil, err
	}
	c.l.Sugar().Debug(zap.Any("res", res))
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
		c.l.Sugar().Error(zap.Error(err))
		return nil, errors.New("error due to unmarshaling json")
	}
	return result, nil
}

func (c *CouriersElasticDAO) Delete(courierID string) error {
	res, err := c.client.Delete().Index(c.index).Type("_doc").Do(context.Background())
	if err != nil {
		c.l.Sugar().Error(err)
		return err
	}
	c.l.Sugar().Debug(zap.Any("res", res))
	return nil
}

func (cd *CouriersElasticDAO) GetMapping() (indexName string, mapping string) {
	return "courier", `{
		"mappings": {
			"_doc": {
				"properties": {
					"name": {
						"type": "keyword"
					},
					"location": {
						"properties": {
							"geo_point": {
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
