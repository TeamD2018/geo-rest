package services

import (
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/TeamD2018/geo-rest/services/photon"
	"github.com/TeamD2018/geo-rest/services/suggestions"
	"go.uber.org/zap"
)

type PhotonSuggestEngineExecutor struct {
	Executors []NamedSuggestEngine
	Logger    *zap.Logger
	Photon    photon.IPhotonClient
}

func NewPhotonSuggestEngineExecutor(photon photon.IPhotonClient, logger *zap.Logger) *PhotonSuggestEngineExecutor {
	return &PhotonSuggestEngineExecutor{
		Photon: photon,
		Logger: logger,
	}
}

func (see *PhotonSuggestEngineExecutor) AddEngine(name string, engine interfaces.SuggestEngine) {
	see.Executors = append(see.Executors, NamedSuggestEngine{name, engine})
}

func (see *PhotonSuggestEngineExecutor) Suggest(input string) (suggestions.SuggestResults, error) {
	results := make(suggestions.SuggestResults)
	for _, executor := range see.Executors {
		req := executor.CreateSearchRequest(input)
		res, err := see.Photon.Search(req.(*photon.SearchQuery))
		if err != nil {
			return nil, err
		}
		if res, err := executor.ParseSearchResponse(res); err != nil {
			return nil, err
		} else {
			results[executor.Name] = res
		}
	}
	return results, nil
}
