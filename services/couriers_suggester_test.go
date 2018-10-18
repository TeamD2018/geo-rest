// +build elastic

package services

import (
	"context"
	"fmt"
	"github.com/TeamD2018/geo-rest/controllers/parameters"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"testing"
)

type CourierSuggesterTestSuite struct {
	suite.Suite
	client       *elastic.Client
	suggester    *CouriersSuggesterElastic
	pool         *dockertest.Pool
	resource     *dockertest.Resource
	logger       *zap.Logger
	testCouriers models.Couriers
	dao          *CouriersElasticDAO
}

var (
	TestPhone812 = "81231231122"
	TestPhone913 = "91331231122"
)

var (
	TestNameABC            = "ABCname"
	TestNameEFG            = "EFGname"
	TestNameWhitespacesABC = "Name ABC"
	TestNameWhitespacesEFG = "Name EFG"
)

var (
	prefs = make(map[string]models.Couriers)
)

func (s *CourierSuggesterTestSuite) BeforeTest(suiteName, testName string) {
	s.suggester.CouriersMapper.EnsureMapping()
	testCourierCreate := []models.CourierCreate{
		{Name: TestNameABC, Phone: &TestPhone812},
		{Name: TestNameABC, Phone: &TestPhone913},
		{Name: TestNameEFG, Phone: &TestPhone812},
		{Name: TestNameEFG, Phone: &TestPhone913},
		{Name: TestNameWhitespacesABC},
		{Name: TestNameWhitespacesEFG},
	}
	s.testCouriers = make(models.Couriers, 0, len(testCourierCreate))
	for _, creation := range testCourierCreate {
		courier, _ := s.dao.Create(&creation)
		s.testCouriers = append(s.testCouriers, courier)
	}
	prefs["ABC"] = models.Couriers{s.testCouriers[0], s.testCouriers[1], s.testCouriers[4]}
	prefs["EFG"] = models.Couriers{s.testCouriers[2], s.testCouriers[3], s.testCouriers[5]}
	prefs["812"] = models.Couriers{s.testCouriers[0], s.testCouriers[2]}
	prefs["913"] = models.Couriers{s.testCouriers[1], s.testCouriers[3]}
	s.client.Refresh(s.dao.GetIndex()).Do(context.Background())
}

func (s *CourierSuggesterTestSuite) AfterTest(suiteName, testName string) {
	s.client.DeleteIndex(s.dao.GetIndex()).Do(context.Background())
}

func (s *CourierSuggesterTestSuite) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func (s *CourierSuggesterTestSuite) SetupSuite() {
	log.SetFlags(log.Lshortfile)
	pool, err := dockertest.NewPool("")
	if err != nil {
		s.FailNow("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("bitnami/elasticsearch", "latest", []string{})
	if err != nil {
		log.Println(err)
		s.FailNow("Could not start resource: %s", err)
	}

	var c *elastic.Client

	if err := pool.Retry(func() error {
		addr := fmt.Sprintf("http://localhost:%s", resource.GetPort("9200/tcp"))

		var err error
		c, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
		if err != nil {
			return err
		}

		_, _, err = c.Ping(addr).Do(context.Background())

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	s.client = c
	s.pool = pool
	s.resource = resource
	s.logger = zap.NewNop()
	s.dao = NewCouriersElasticDAO(s.client, s.logger, "", 0)
	s.suggester = NewCouriersSuggesterElastic(s.client, s.dao, s.logger)
	s.suggester.SetFuzziness(1)
}

func (s *CourierSuggesterTestSuite) TestCourierSuggesterTestSuite_Suggest_ByPhone_OK() {
	service := s.suggester

	got, err := service.Suggest("suggestions", &parameters.Suggestion{Prefix: "913", Limit: 200})
	s.NoError(err)
	want := prefs["913"]
	s.ElementsMatch(got, want)

	s.logger.Debug("arrays", zap.Any("got", got), zap.Any("want", want))
}

func (s *CourierSuggesterTestSuite) TestCourierSuggesterTestSuite_Suggest_ByName_OK() {
	service := s.suggester
	got, err := service.Suggest("suggestions", &parameters.Suggestion{Prefix: "ABC", Limit: 200})
	s.NoError(err)
	want := prefs["ABC"]
	s.ElementsMatch(got, want)

	s.logger.Debug("arrays", zap.Any("got", got), zap.Any("want", want))
}

func (s *CourierSuggesterTestSuite) TestCourierSuggesterTestSuite_Suggest_ByNameLowerCase_OK() {
	service := s.suggester

	got, err := service.Suggest("suggestions", &parameters.Suggestion{Prefix: "abc", Limit: 200})
	s.NoError(err)
	want := prefs["ABC"]

	s.logger.Debug("arrays", zap.Any("got", got), zap.Any("want", want))
	s.ElementsMatch(got, want)

}

func (s *CourierSuggesterTestSuite) TestCourierSuggesterTestSuite_Suggest_ByNameFuzzy_OK() {
	service := s.suggester
	got, err := service.Suggest("suggestions", &parameters.Suggestion{Prefix: "ADC", Limit: 200})
	s.NoError(err)
	want := prefs["ABC"]
	s.ElementsMatch(got, want)
}

func (s *CourierSuggesterTestSuite) TestCourierSuggesterTestSuite_Suggest_ByPhoneFuzzy_OK() {
	service := s.suggester
	got, err := service.Suggest("suggestions", &parameters.Suggestion{Prefix: "9134", Limit: 200})
	s.NoError(err)
	want := prefs["913"]
	s.ElementsMatch(got, want)

	s.logger.Debug("arrays", zap.Any("got", got), zap.Any("want", want))
}

func TestIntegrationSuggesterSuite(t *testing.T) {
	suite.Run(t, new(CourierSuggesterTestSuite))
}
