package interfaces

import (
	"encoding/json"
	"github.com/olivere/elastic"
)

type SuggestOptions struct {
}

type SuggestResult struct {
	Source *json.RawMessage
	Id     string
}

type EngineSuggestResults []SuggestResult
type SuggestResults map[string]EngineSuggestResults

type SuggestEngine interface {
	CreateSearchRequest(input string) (*elastic.SearchRequest)
	ParseSearchResponse(result *elastic.SearchResult) EngineSuggestResults
}

type SuggestExecutor interface {
	AddEngine(name string, engine SuggestEngine)
	Suggest(input string) (SuggestResults, error)
}
