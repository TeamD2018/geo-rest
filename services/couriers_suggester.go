package services

import (
	"context"
	"encoding/json"
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
)

const Suggester string = "couriers-suggester"

type CouriersSuggesterElastic struct {
	Elastic        *elastic.Client
	CouriersMapper interfaces.Mapper
	logger         *zap.Logger
	fuzziness      int
}

func NewCouriersSuggesterElastic(client *elastic.Client, mapper interfaces.Mapper, logger *zap.Logger) *CouriersSuggesterElastic {
	return &CouriersSuggesterElastic{
		Elastic:        client,
		CouriersMapper: mapper,
		logger:         logger,
		fuzziness:      0,
	}
}

func (cs *CouriersSuggesterElastic) Suggest(field string, suggestion *parameters.Suggestion) (models.Couriers, error) {
	db := cs.Elastic
	suggester := elastic.NewCompletionSuggester(Suggester).
		Field(field).
		PrefixWithEditDistance(suggestion.Prefix, cs.fuzziness).
		Size(suggestion.Limit)

	query := db.Search(cs.CouriersMapper.GetIndex()).Type("_doc").Suggester(suggester)
	res, err := query.Do(context.Background())
	if err != nil {
		if elastic.IsNotFound(err) {
			return make(models.Couriers, 0), nil
		}
		cs.logger.Error("fail to suggest", zap.Error(err))
		return nil, err
	}
	suggestions := res.Suggest[Suggester]
	found := make(models.Couriers, 0)
	for _, suggest := range suggestions {
		for _, option := range suggest.Options {
			var courier models.Courier
			if err := json.Unmarshal(*option.Source, &courier); err != nil {
				return nil, err
			}
			courier.ID = option.Id
			found = append(found, &courier)
		}
	}
	return found, err
}

func (cs *CouriersSuggesterElastic) SetFuzziness(fuzziness int) {
	cs.fuzziness = fuzziness
}
