package interfaces

import (
	"github.com/TeamD2018/geo-rest/services/suggestions"
	"github.com/olivere/elastic"
)

type SuggestEngine interface {
	CreateSearchRequest(input string) (*elastic.SearchRequest)
	ParseSearchResponse(result *elastic.SearchResult) suggestions.EngineSuggestResults
}

type SuggestExecutor interface {
	AddEngine(name string, engine SuggestEngine)
	Suggest(input string) (suggestions.SuggestResults, error)
}
