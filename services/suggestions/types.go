package suggestions

import "encoding/json"

type ElasticSuggestResult struct {
	Source *json.RawMessage
	Id     string
}

type SuggestResults map[string]interface{}
