package services

import (
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/TeamD2018/geo-rest/services/suggestions"
)

type SuggestionService struct {
	Executors []interfaces.SuggestExecutor
}

func NewSuggestionService(executors ... interfaces.SuggestExecutor) *SuggestionService {
	return &SuggestionService{
		Executors: executors,
	}
}

func (ss *SuggestionService) Suggest(input string) (suggestions.SuggestResults, error) {
	results := make(suggestions.SuggestResults)
	for _, executor := range ss.Executors {
		suggestion, err := executor.Suggest(input)
		if err != nil {
			return nil, err
		}
		for k, v := range suggestion {
			results[k] = v
		}
	}
	return results, nil
}
