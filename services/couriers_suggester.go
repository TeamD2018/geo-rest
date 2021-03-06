package services

import (
	"context"
	"encoding/json"
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"strings"
)

const Suggester string = "couriers-suggester"

type CouriersSuggesterElastic struct {
	Elastic        *elastic.Client
	CouriersMapper interfaces.Mapper
	logger         *zap.Logger
	fuzziness      int
	threshold      int
}

const CouriersDefaultFuzziness = 2
const CouriersDefaultFuzzinessThreshold = 5

func (cs *CouriersSuggesterElastic) SetFuzzinessThreshold(threshold int) interfaces.CourierSuggester {
	cs.threshold = threshold
	return cs
}

func NewCouriersSuggesterElastic(client *elastic.Client, mapper interfaces.Mapper, logger *zap.Logger) *CouriersSuggesterElastic {
	return &CouriersSuggesterElastic{
		Elastic:        client,
		CouriersMapper: mapper,
		logger:         logger,
		fuzziness:      CouriersDefaultFuzziness,
		threshold:      CouriersDefaultFuzzinessThreshold,
	}
}

func (cs *CouriersSuggesterElastic) Suggest(field string, suggestion *parameters.Suggestion) (models.Couriers, error) {
	db := cs.Elastic
	fuzzyOptions := elastic.NewFuzzyCompletionSuggesterOptions().
		MinLength(cs.threshold).
		EditDistance(cs.fuzziness).
		UnicodeAware(true)
	suggester := elastic.NewCompletionSuggester(Suggester).
		Field(field).
		FuzzyOptions(fuzzyOptions).
		Prefix(strings.ToLower(suggestion.Prefix)).
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

func (cs *CouriersSuggesterElastic) SetFuzziness(fuzziness int) interfaces.CourierSuggester {
	cs.fuzziness = fuzziness
	return cs
}
