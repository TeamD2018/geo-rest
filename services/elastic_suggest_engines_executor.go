package services

import (
	"context"
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/TeamD2018/geo-rest/services/suggestions"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
)

type ElasticSuggestEngineExecutor struct {
	Executors []NamedSuggestEngine
	Logger    *zap.Logger
	Elastic   *elastic.Client
}

type NamedSuggestEngine struct {
	Name string
	interfaces.SuggestEngine
}

func NewSuggestEngineExecutor(client *elastic.Client, logger *zap.Logger) *ElasticSuggestEngineExecutor {
	return &ElasticSuggestEngineExecutor{
		Elastic: client,
		Logger:  logger,
	}
}

func (see *ElasticSuggestEngineExecutor) AddEngine(name string, engine interfaces.SuggestEngine) {
	see.Executors = append(see.Executors, NamedSuggestEngine{name, engine})
}

func (see *ElasticSuggestEngineExecutor) Suggest(input string) (suggestions.SuggestResults, error) {
	results := make(suggestions.SuggestResults)
	multisearch := see.Elastic.MultiSearch()
	total := 0
	for _, executor := range see.Executors {
		req := executor.CreateSearchRequest(input)
		multisearch.Add(req.(*elastic.SearchRequest))
		total++
	}
	res, err := multisearch.MaxConcurrentSearches(total).Do(context.Background())

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
