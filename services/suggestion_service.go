package services

import (
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/TeamD2018/geo-rest/services/suggestions"
	"sync"
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
	resChan := make(chan suggestions.SuggestResults, len(ss.Executors))
	errc := make(chan error, 1)
	var wg sync.WaitGroup
	for _, executor := range ss.Executors {
		wg.Add(1)
		go func(
			in string,
			ex interfaces.SuggestExecutor,
			res chan<- suggestions.SuggestResults,
			e chan<- error) {
			defer wg.Done()
			ss.runExecutor(in, ex, res, e)
		}(input, executor, resChan, errc)
	}
	wg.Wait()
	close(errc)
	close(resChan)
	err, more := <-errc
	if more {
		return nil, err
	}
	results := make(suggestions.SuggestResults)
	for result := range resChan {
		for k, v := range result {
			results[k] = v
		}
	}
	return results, nil
}

func (ss *SuggestionService) runExecutor(
	input string,
	executor interfaces.SuggestExecutor,
	result chan<- suggestions.SuggestResults,
	errc chan<- error) {
	suggestion, err := executor.Suggest(input)
	if err != nil {
		select {
		case errc <- err:
			return
		default:
			return
		}
	}
	result <- suggestion
}
