package interfaces

import (
	"github.com/TeamD2018/geo-rest/models"
	"github.com/TeamD2018/geo-rest/services/suggestions"
)

type SuggestEngine interface {
	CreateSearchRequest(input string) (interface{})
	ParseSearchResponse(result interface{}) (interface{}, error)
}

type SuggestExecutor interface {
	AddEngine(name string, engine SuggestEngine)
	SuggestionService
}

type SuggestionService interface {
	Suggest(input string) (suggestions.SuggestResults, error)
}

type IConcurrentLookupService interface {
	LookupAll(
		source <-chan *models.OSMEntity,
		destination chan<- *models.OSMPolygonSuggestion,
		errc chan<- error) func()
}
