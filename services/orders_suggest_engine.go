package services

import (
	"github.com/TeamD2018/geo-rest/services/suggestions"
	"github.com/olivere/elastic"
)

type OrdersSuggestEngine struct {
	Field              string
	Fuzziness          string
	Index              string
	Limit              int
	FuzzinessThreshold int
}

func (ose *OrdersSuggestEngine) ParseSearchResponse(response interface{}) (interface{}, error) {
	result := response.(*elastic.SearchResult)
	if result.TotalHits() == 0 {
		return nil, nil
	}
	results := make([]suggestions.ElasticSuggestResult, 0, result.TotalHits())
	for _, hit := range result.Hits.Hits {
		results = append(results, suggestions.ElasticSuggestResult{Id: hit.Id, Source: hit.Source})
	}
	return results, nil
}

func (ose *OrdersSuggestEngine) CreateSearchRequest(input string) (interface{}) {
	query := elastic.NewMatchQuery(ose.Field, input).Operator("and")
	if len(input) >= ose.FuzzinessThreshold {
		query.Fuzziness(ose.Fuzziness)
	}
	source := elastic.NewSearchSource().Query(query).Size(ose.Limit)
	return elastic.NewSearchRequest().SearchSource(source).Index(ose.Index).Type("_doc")
}
