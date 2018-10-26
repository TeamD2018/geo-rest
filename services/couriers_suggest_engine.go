package services

import (
	"github.com/TeamD2018/geo-rest/services/interfaces"
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
		MinLength(cse.FuzzinessThreshold).
		EditDistance(cse.Fuzziness).
		UnicodeAware(true)

	suggester := elastic.NewCompletionSuggester(Suggester).
		Field(cse.Field).
		FuzzyOptions(fuzzyOptions).
		Prefix(strings.ToLower(input)).
		Size(cse.Limit)
	source := elastic.NewSearchSource().Suggester(suggester)
	return elastic.NewSearchRequest().SearchSource(source).Index(cse.Index).Type("_doc")
}

func (cse *CouriersSuggestEngine) ParseSearchResponse(result *elastic.SearchResult) interfaces.EngineSuggestResults {
	suggestResults := result.Suggest[CouriersSuggesterName]
	results := make(interfaces.EngineSuggestResults, 0)
	for _, suggestion := range suggestResults {
		for _, option := range suggestion.Options {
			results = append(results, interfaces.SuggestResult{Id: option.Id, Source: option.Source})
		}
	}
	return results
}
