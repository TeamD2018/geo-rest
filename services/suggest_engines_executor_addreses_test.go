package services

import (
	"context"
	"fmt"
	"github.com/TeamD2018/geo-rest/models"
	"github.com/olivere/elastic"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"testing"
)

type TestSuggestEnginesExecutor struct {
	suite.Suite
	client         *elastic.Client
	ordersDao      *OrdersElasticDAO
	couriersDao    *CouriersElasticDAO
	ordersEngine   OrdersSuggestEngine
	couriersEngine CouriersSuggestEngine
	executor       *SuggestEngineExecutor
	pool           *dockertest.Pool
	resource       *dockertest.Resource
	logger         *zap.Logger
}

const (
	ordersEngine   = "orders-engine"
	couriersEngine = "couriers-engine"
)

var (
	courierId         string
	testCourierCreate = models.CourierCreate{Name: "test-courier"}
	testCourier       *models.Courier
)

var (
	TestAddressSlash                     = "ул. Николая Химушина, 2/7 строение 1, Москва, 107143"
	TestAddressNoStreetKeyword           = "1 этаж, отдел 41, пр. Вернадского, 14А, Москва, 119415"
	TestAddressSimpleWithHouse           = "Тагильская улица, дом 2, Москва, 107143"
	TestAddressWithDistrict              = "Северная ул., 12, Ромашково, Московская обл., Россия, 143025"
	TestAddressWithSameDistrict          = "Тагильская улица, 7, Ромашково, Московская обл., Россия, 143025"
	TestAddressComplexWithHouseShorthand = "Большая Черкизовская ул., д93А, Москва, 107553"
	TestAddressComplexWithHouse          = "Большая Черкизовская улица, дом 94А, Москва, 107553"
)

var (
	AddressSlashInput                = "Химушина 2/7"
	AddressSlashInputReversed        = "2/7 Химушина"
	AddressSlashInputWithStreet      = "ул Химушина 2/7"
	AddressSlashInputWithStreetPunct = "ул. Химушина 2/7"

	AddressNoStreetKeywordInput         = "Вернадского 14А"
	AddressNoStreetKeywordInputReversed = "14А Вернадского"
	AddressNoStreetKeywordAdditional    = "отдел 41 Вернадского"
	AddressNoStreetWithAvenueShorthand  = "пр Вернадского 14А"
	AddressNoStreetWithAvenueFull       = "проспект Вернадского 14А"

	AddressSimpleHouseInput                      = "Тагильская дом 2"
	AddressSimpleHouseMinimalInput               = "Тагильская 2"
	AddressSimpleShorthandHouseInput             = "Тагильская д2"
	AddressStreetShorthandInput                  = "у Тагильская"
	AddressStreetAndHouseShorthandInput          = "у Тагильская д2"
	AddressStreetAndHouseShorthandInputWithSpace = "у Тагильская д 2"

	AddressWithDistrictInput              = "Ромашково"
	AddressComplexWithHouseInput          = "94А Черкизовская"
	AddressComplexWithHouseLowercaseInput = "94а Черкизовская"
	AddressComplexWithHouseShorthandInput = "д93А Черкизовская"
)

var (
	ordersOutputs = make(map[string]models.Orders)
)

var (
	orderAddressSlash                     *models.Order
	orderAddressNoStreetKeyword           *models.Order
	orderAddressSimpleWithHouse           *models.Order
	orderAddressWithDistrict              *models.Order
	orderAddressWithSameDistrict          *models.Order
	orderAddressWithHouse                 *models.Order
	orderAddressComplexWithHouseShorthand *models.Order
)

func (s *TestSuggestEnginesExecutor) SetupSuite() {
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
	s.couriersDao = NewCouriersElasticDAO(s.client, s.logger, "couriers", 100)
	testCourier, _ = s.couriersDao.Create(&testCourierCreate)
	courierId = testCourier.ID
	s.ordersDao = NewOrdersElasticDAO(s.client, s.logger, s.couriersDao, "orders")
	s.ordersEngine = OrdersSuggestEngine{
		Field:              "destination.address",
		Fuzziness:          "0",
		Limit:              200,
		Index:              s.ordersDao.GetIndex(),
		FuzzinessThreshold: 3,
	}
	s.executor = NewSuggestEngineExecutor(s.client, s.logger)
	s.executor.AddEngine(ordersEngine, &s.ordersEngine)
	var (
		createOrderAddressSlash                     = &models.OrderCreate{CourierID: &courierId, Destination: models.Location{Address: &TestAddressSlash}}
		createOrderAddressNoStreetKeyword           = &models.OrderCreate{CourierID: &courierId, Destination: models.Location{Address: &TestAddressNoStreetKeyword}}
		createOrderAddressSimpleWithHouse           = &models.OrderCreate{CourierID: &courierId, Destination: models.Location{Address: &TestAddressSimpleWithHouse}}
		createOrderAddressWithDistrict              = &models.OrderCreate{CourierID: &courierId, Destination: models.Location{Address: &TestAddressWithDistrict}}
		createOrderAddressWithSameDistrict          = &models.OrderCreate{CourierID: &courierId, Destination: models.Location{Address: &TestAddressWithSameDistrict}}
		createOrderAddressComplexWithHouseShorthand = &models.OrderCreate{CourierID: &courierId, Destination: models.Location{Address: &TestAddressComplexWithHouseShorthand}}
		createOrderAddressComplexWithHouse          = &models.OrderCreate{CourierID: &courierId, Destination: models.Location{Address: &TestAddressComplexWithHouse}}
	)
	if err := s.ordersDao.EnsureMapping(); err != nil {
		s.FailNow("orders dao  mapping failed", err.Error())
	}
	orderAddressSlash, _ := s.ordersDao.Create(createOrderAddressSlash)
	orderAddressNoStreetKeyword, _ = s.ordersDao.Create(createOrderAddressNoStreetKeyword)
	orderAddressSimpleWithHouse, _ = s.ordersDao.Create(createOrderAddressSimpleWithHouse)
	orderAddressWithDistrict, _ = s.ordersDao.Create(createOrderAddressWithDistrict)
	orderAddressWithSameDistrict, _ = s.ordersDao.Create(createOrderAddressWithSameDistrict)
	orderAddressWithHouse, _ = s.ordersDao.Create(createOrderAddressComplexWithHouse)
	orderAddressComplexWithHouseShorthand, _ = s.ordersDao.Create(createOrderAddressComplexWithHouseShorthand)

	s.ordersDao.Elastic.Refresh("_all").Do(context.Background())
	ordersOutputs = map[string]models.Orders{
		"AddressSlashInput":                {orderAddressSlash},
		"AddressSlashInputReversed":        {orderAddressSlash},
		"AddressSlashInputWithStreet":      {orderAddressSlash},
		"AddressSlashInputWithStreetPunct": {orderAddressSlash},

		"AddressNoStreetKeywordInput":         {orderAddressNoStreetKeyword},
		"AddressNoStreetKeywordInputReversed": {orderAddressNoStreetKeyword},
		"AddressNoStreetKeywordAdditional":    {orderAddressNoStreetKeyword},
		"AddressNoStreetWithAvenueShorthand":  {orderAddressNoStreetKeyword},
		"AddressNoStreetWithAvenueFull":       {orderAddressNoStreetKeyword},

		"AddressSimpleHouseInput":                      {orderAddressSimpleWithHouse},
		"AddressSimpleShorthandHouseInput":             {orderAddressSimpleWithHouse},
		"AddressStreetShorthandInput":                  {orderAddressSimpleWithHouse, orderAddressWithSameDistrict},
		"AddressStreetAndHouseShorthandInput":          {orderAddressSimpleWithHouse},
		"AddressStreetAndHouseShorthandInputWithSpace": {orderAddressSimpleWithHouse},

		"AddressWithDistrictInput":              {orderAddressWithDistrict, orderAddressWithSameDistrict},
		"AddressComplexWithHouseInput":          {orderAddressWithHouse},
		"AddressComplexWithHouseLowercaseInput": {orderAddressWithHouse},
		"AddressComplexWithHouseShorthandInput": {orderAddressComplexWithHouseShorthand},

		"AddressSimpleHouseMinimalInput": {orderAddressSimpleWithHouse},
	}

}

func (s *TestSuggestEnginesExecutor) TearDownSuite() {
	s.Nil(s.pool.Purge(s.resource))
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressSlash() {
	results, err := s.executor.Suggest(AddressSlashInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressSlashInput"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressSlashReversed() {
	results, err := s.executor.Suggest(AddressSlashInputReversed)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressSlashInputReversed"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressSlashInputWithStreet() {
	results, err := s.executor.Suggest(AddressSlashInputWithStreet)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressSlashInputWithStreet"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressSlashInputWithStreetPunct() {
	results, err := s.executor.Suggest(AddressSlashInputWithStreetPunct)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressSlashInputWithStreetPunct"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressNoStreetKeywordInput() {
	results, err := s.executor.Suggest(AddressNoStreetKeywordInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressNoStreetKeywordInput"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressNoStreetKeywordInputReversed() {
	results, err := s.executor.Suggest(AddressNoStreetKeywordInputReversed)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressNoStreetKeywordInputReversed"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressNoStreetKeywordAdditional() {
	results, err := s.executor.Suggest(AddressNoStreetKeywordAdditional)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressNoStreetKeywordAdditional"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressNoStreetWithAvenueShorthand() {
	results, err := s.executor.Suggest(AddressNoStreetWithAvenueShorthand)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressNoStreetWithAvenueShorthand"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressNoStreetWithAvenueFull() {
	results, err := s.executor.Suggest(AddressNoStreetWithAvenueFull)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressNoStreetWithAvenueFull"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressSimpleHouseInput() {
	results, err := s.executor.Suggest(AddressSimpleHouseInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressSimpleHouseInput"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressSimpleShorhandHouseINput() {
	results, err := s.executor.Suggest(AddressSimpleShorthandHouseInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressSimpleShorthandHouseInput"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressStreetShorthandInput() {
	results, err := s.executor.Suggest(AddressStreetShorthandInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressStreetShorthandInput"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressStreetAndHouseShorthandInput() {
	results, err := s.executor.Suggest(AddressStreetAndHouseShorthandInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressStreetAndHouseShorthandInput"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressStreetAndHouseShorthandInputWithSpace() {
	results, err := s.executor.Suggest(AddressStreetAndHouseShorthandInputWithSpace)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressStreetAndHouseShorthandInputWithSpace"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressWithDistrictInput() {
	results, err := s.executor.Suggest(AddressWithDistrictInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressWithDistrictInput"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressComplexWithHouseInput() {
	results, err := s.executor.Suggest(AddressComplexWithHouseInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressComplexWithHouseInput"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressComplexWithHouseLowercaseInput() {
	results, err := s.executor.Suggest(AddressComplexWithHouseLowercaseInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressComplexWithHouseLowercaseInput"]
	s.ElementsMatch(want, got.Orders)
}

func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressComplexWithHouseShorthandInput() {
	results, err := s.executor.Suggest(AddressComplexWithHouseShorthandInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressComplexWithHouseShorthandInput"]
	s.ElementsMatch(want, got.Orders)
}
func (s *TestSuggestEnginesExecutor) TestSuggestEnginesExecutor_Suggest_TestAddressSimpleHouseMinimalInput() {
	results, err := s.executor.Suggest(AddressSimpleHouseMinimalInput)
	if !s.NoError(err) {
		return
	}
	orders := results[ordersEngine]
	couriers := results[couriersEngine]
	got, err := models.SuggestionFromRawInput(orders, couriers)
	if !s.NoError(err) {
		return
	}
	want := ordersOutputs["AddressSimpleHouseMinimalInput"]
	s.ElementsMatch(want, got.Orders)
}
func TestIntegrationSuggestEnginesExecutor(t *testing.T) {
	suite.Run(t, new(TestSuggestEnginesExecutor))
}
