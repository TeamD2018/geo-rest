package services

import (
	"github.com/TeamD2018/geo-rest/services/interfaces"
	"github.com/TeamD2018/geo-rest/services/suggestions"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuggestionServiceTestSuite struct {
	suite.Suite
	service *SuggestionService
}

type engineExecutorMock struct {
	mock.Mock
}

func (em *engineExecutorMock) AddEngine(name string, engine interfaces.SuggestEngine) {
	em.Called(name, engine)
	return
}

func (em engineExecutorMock) Suggest(input string) (suggestions.SuggestResults, error) {
	args := em.Called(input)
	return args.Get(0).(suggestions.SuggestResults), args.Error(1)
}

func (s *SuggestionServiceTestSuite) SetupMock(expected suggestions.SuggestResults) {
	var mocks []interfaces.SuggestExecutor
	for k, v := range expected {
		executor := new(engineExecutorMock)
		executor.On("Suggest", mock.AnythingOfType("string")).Return(suggestions.SuggestResults{k: v}, nil)
		mocks = append(mocks, executor)
	}
	s.service = NewSuggestionService(mocks...)
}

func (s *SuggestionServiceTestSuite) TestSuggestionsService_Suggest_OK() {
	expected := suggestions.SuggestResults{
		"expected_1": "value_1",
		"expected_2": "value_2",
		"expected_3": "value_3",
	}
	s.SetupMock(expected)
	res, err := s.service.Suggest("suggest")
	if !s.NoError(err) {
		return
	}
	s.EqualValues(expected, res)
}

func TestSuggestionService(t *testing.T) {
	suite.Run(t, new(SuggestionServiceTestSuite))
}
