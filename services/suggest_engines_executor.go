package services

import (
	"context"
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/olivere/elastic"
)

type SuggestEngineExecutor struct {
	Executors []NamedSuggestEngine
	Elastic   *elastic.Client
}

type NamedSuggestEngine struct {
	Name string
	interfaces.SuggestEngine
}

func (see *SuggestEngineExecutor) AddEngine(name string, engine interfaces.SuggestEngine) {
	see.Executors = append(see.Executors, NamedSuggestEngine{name, engine})
}

func (see *SuggestEngineExecutor) Suggest(input string) (interfaces.SuggestResults, error) {
	results := make(interfaces.SuggestResults)
	mutlisearch := see.Elastic.MultiSearch()
	total := 0
	for _, executor := range see.Executors {
		mutlisearch.Add(executor.CreateSearchRequest(input))
		total++
	}
	res, err := mutlisearch.MaxConcurrentSearches(total).Do(context.Background())
	if err != nil {
		return nil, err
	}
	for i, response := range res.Responses {
		executor := see.Executors[i]
		executorName := executor.Name
		results[executorName] = executor.ParseSearchResponse(response)
	}
	return results, nil
}
