package suggestions

import "encoding/json"

type SuggestResult struct {
	Source *json.RawMessage
	Id     string
}

type EngineSuggestResults []SuggestResult
type SuggestResults map[string]EngineSuggestResults
