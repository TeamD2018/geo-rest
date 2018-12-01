package interfaces

import (
	"github.com/TeamD2018/geo-rest/services/suggestions"
)

type SuggestEngine interface {
	CreateSearchRequest(input string) (interface{})
	ParseSearchResponse(result interface{}) interface{}
}

type SuggestExecutor interface {
	AddEngine(name string, engine SuggestEngine)
	SuggestionService
}

type SuggestionService interface {
	Suggest(input string) (suggestions.SuggestResults, error)
}
