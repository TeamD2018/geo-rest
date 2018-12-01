package services

import (
	"github.com/TeamD2018/geo-rest/services/suggestions"
	"github.com/olivere/elastic"
	"strings"
)

type PrefixSuggestEngine struct {
	Field              string
	Fuzziness          string
	Index              string
	Limit              int
	FuzzinessThreshold int
}

func (ops *PrefixSuggestEngine) ParseSearchResponse(response interface{}) interface{} {
	result := response.(*elastic.SearchResult)
	suggestResults := result.Suggest[CouriersSuggesterName]
	results := make([]suggestions.ElasticSuggestResult, 0)
	for _, suggestion := range suggestResults {
		for _, option := range suggestion.Options {
			results = append(results, suggestions.ElasticSuggestResult{Id: option.Id, Source: option.Source})
		}
	}
	return results
}

func (ops *PrefixSuggestEngine) CreateSearchRequest(input string) (interface{}) {
	fuzzyOptions := elastic.NewFuzzyCompletionSuggesterOptions().
		UnicodeAware(true)
	if ops.FuzzinessThreshold > 0 {
		fuzzyOptions.MinLength(ops.FuzzinessThreshold)
	}
	if ops.Fuzziness != "" {
		fuzzyOptions.EditDistance(ops.Fuzziness)
	}

	suggester := elastic.NewCompletionSuggester(CouriersSuggesterName).
		Field(ops.Field).
		FuzzyOptions(fuzzyOptions).
		Prefix(strings.ToLower(input)).
		Size(ops.Limit)
	source := elastic.NewSearchSource().Suggester(suggester)
	return elastic.NewSearchRequest().SearchSource(source).Index(ops.Index).Type("_doc")
}
