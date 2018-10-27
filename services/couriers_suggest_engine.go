package services

import (
	"github.com/TeamD2018/geo-rest/services/suggestions"
	"github.com/olivere/elastic"
	"strings"
)

const CouriersSuggesterName = "couriers-suggest"

type CouriersSuggestEngine struct {
	Field              string
	Fuzziness          string
	FuzzinessThreshold int
	Limit              int
	Index              string
}

func (cse *CouriersSuggestEngine) CreateSearchRequest(input string) (*elastic.SearchRequest) {
	fuzzyOptions := elastic.NewFuzzyCompletionSuggesterOptions().
		UnicodeAware(true)
	if cse.FuzzinessThreshold > 0 {
		fuzzyOptions.MinLength(cse.FuzzinessThreshold)
	}
	if cse.Fuzziness != "" {
		fuzzyOptions.EditDistance(cse.Fuzziness)
	}

	suggester := elastic.NewCompletionSuggester(CouriersSuggesterName).
		Field(cse.Field).
		FuzzyOptions(fuzzyOptions).
		Prefix(strings.ToLower(input)).
		Size(cse.Limit)
	source := elastic.NewSearchSource().Suggester(suggester)
	return elastic.NewSearchRequest().SearchSource(source).Index(cse.Index).Type("_doc")
}

func (cse *CouriersSuggestEngine) ParseSearchResponse(result *elastic.SearchResult) suggestions.EngineSuggestResults {
	suggestResults := result.Suggest[CouriersSuggesterName]
	results := make(suggestions.EngineSuggestResults, 0)
	for _, suggestion := range suggestResults {
		for _, option := range suggestion.Options {
			results = append(results, suggestions.SuggestResult{Id: option.Id, Source: option.Source})
		}
	}
	return results
}
