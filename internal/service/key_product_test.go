package service

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	casbinMocks "github.com/paysuper/casbin-server/pkg/mocks"
	"github.com/paysuper/paysuper-billing-server/internal/config"
	"github.com/paysuper/paysuper-billing-server/internal/database"
	"github.com/paysuper/paysuper-billing-server/internal/mocks"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	reportingMocks "github.com/paysuper/paysuper-proto/go/reporterpb/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"gopkg.in/ProtocolONE/rabbitmq.v1/pkg"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"testing"
)

type KeyProductTestSuite struct {
	suite.Suite
	service *Service
	log     *zap.Logger
	cache   database.CacheInterface

	project    *billingpb.Project
	pmBankCard *billingpb.PaymentMethod
}

func Test_KeyProduct(t *testing.T) {
	suite.Run(t, new(KeyProductTestSuite))
}

func (suite *KeyProductTestSuite) SetupTest() {
	cfg, err := config.NewConfig()
	assert.NoError(suite.T(), err, "Config load failed")

	db, err := mongodb.NewDatabase()
	if err != nil {
		suite.FailNow("Database connection failed", "%v", err)
	}

	m, err := migrate.New(
		"file://../../migrations/tests/keys",
		cfg.MongoDsn)
	assert.NoError(suite.T(), err, "Migrate init failed")
	if err != nil {
		suite.FailNow("Migrations failed", "%v", err)
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		suite.FailNow("Migrations failed", "%v", err)
	}

	pgRub := &billingpb.PriceGroup{
		Id:       primitive.NewObjectID().Hex(),
		Region:   "RUB",
		Currency: "RUB",
		IsActive: true,
	}
	pgUsd := &billingpb.PriceGroup{
		Id:       primitive.NewObjectID().Hex(),
		Region:   "USD",
		Currency: "USD",
		IsActive: true,
	}
	pgEur := &billingpb.PriceGroup{
		Id:       primitive.NewObjectID().Hex(),
		Region:   "EUR",
		Currency: "EUR",
		IsActive: true,
	}
	if err != nil {
		suite.FailNow("Insert currency test data failed", "%v", err)
	}

	suite.log, err = zap.NewProduction()
	assert.NoError(suite.T(), err, "Logger initialization failed")

	broker, err := rabbitmq.NewBroker(cfg.BrokerAddress)
	assert.NoError(suite.T(), err, "Creating RabbitMQ publisher failed")

	redisdb := mocks.NewTestRedis()
	suite.cache, err = database.NewCacheRedis(redisdb, "cache")
	suite.service = NewBillingService(
		db,
		cfg,
		mocks.NewGeoIpServiceTestOk(),
		mocks.NewRepositoryServiceOk(),
		mocks.NewTaxServiceOkMock(),
		broker,
		nil,
		suite.cache,
		mocks.NewCurrencyServiceMockOk(),
		mocks.NewDocumentSignerMockOk(),
		&reportingMocks.ReporterService{},
		mocks.NewFormatterOK(),
		mocks.NewBrokerMockOk(),
		&casbinMocks.CasbinService{},
	)

	if err := suite.service.Init(); err != nil {
		suite.FailNow("Billing service initialization failed", "%v", err)
	}

	suite.NoError(suite.service.merchant.Insert(ctx, &billingpb.Merchant{Id: merchantId, Banking: &billingpb.MerchantBanking{Currency: "USD"}}))

	pgs := []*billingpb.PriceGroup{pgRub, pgUsd, pgEur}
	if err := suite.service.priceGroupRepository.MultipleInsert(ctx, pgs); err != nil {
		suite.FailNow("Insert price group test data failed", "%v", err)
	}

	if err := suite.service.project.Insert(context.TODO(), &billingpb.Project{
		Id:         projectId,
		MerchantId: merchantId,
	}); err != nil {
		suite.FailNow("Insert project test data failed", "%v", err)
	}
}

func (suite *KeyProductTestSuite) TearDownTest() {
	err := suite.service.db.Drop()

	if err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	err = suite.service.db.Close()

	if err != nil {
		suite.FailNow("Database close failed", "%v", err)
	}
}

func (suite *KeyProductTestSuite) Test_GetKeyProductInfo() {
	shouldBe := require.New(suite.T())

	req := &billingpb.CreateOrUpdateKeyProductRequest{
		Object:          "product",
		Sku:             "ru_double_yeti",
		Name:            map[string]string{"en": initialName},
		DefaultCurrency: "USD",
		Description:     map[string]string{"en": "blah-blah-blah"},
		LongDescription: map[string]string{"en": "Super game steam keys"},
		Url:             "http://test.ru/dffdsfsfs",
		Cover: &billingpb.ImageCollection{
			UseOneForAll: false,
			Images: &billingpb.LocalizedUrl{
				En: "/home/image.jpg",
			},
		},
		MerchantId: merchantId,
		ProjectId:  projectId,
		Platforms: []*billingpb.PlatformPrice{
			{
				Id: "steam",
				Prices: []*billingpb.ProductPrice{
					{Region: "USD", Currency: "USD", Amount: 10},
					{Region: "EUR", Currency: "EUR", Amount: 20},
				},
			},
		},
		Metadata: map[string]string{
			"SomeKey": "SomeValue",
		},
	}
	response := billingpb.KeyProductResponse{}
	err := suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, &response)
	shouldBe.Nil(err)
	shouldBe.Nil(response.Message)

	res := billingpb.GetKeyProductInfoResponse{}
	err = suite.service.GetKeyProductInfo(context.TODO(), &billingpb.GetKeyProductInfoRequest{Currency: "USD", KeyProductId: response.Product.Id, Language: "en"}, &res)
	shouldBe.Nil(err)
	shouldBe.NotNil(res.Message)
	shouldBe.EqualValues(400, res.Status)

	publishRsp := &billingpb.KeyProductResponse{}
	err = suite.service.PublishKeyProduct(context.TODO(), &billingpb.PublishKeyProductRequest{MerchantId: merchantId, KeyProductId: response.Product.Id}, publishRsp)
	shouldBe.Nil(err)
	shouldBe.EqualValues(200, publishRsp.Status)

	res = billingpb.GetKeyProductInfoResponse{}
	err = suite.service.GetKeyProductInfo(context.TODO(), &billingpb.GetKeyProductInfoRequest{Currency: "USD", KeyProductId: response.Product.Id, Language: "en"}, &res)
	shouldBe.Nil(err)
	shouldBe.Nil(res.Message)
	shouldBe.NotNil(res.KeyProduct)
	shouldBe.Equal(response.Product.Id, res.KeyProduct.Id)
	shouldBe.Equal(initialName, res.KeyProduct.Name)
	shouldBe.Equal("blah-blah-blah", res.KeyProduct.Description)
	shouldBe.Equal(1, len(res.KeyProduct.Platforms))
	shouldBe.Equal("steam", res.KeyProduct.Platforms[0].Id)
	shouldBe.EqualValues(10, res.KeyProduct.Platforms[0].Price.Amount)
	shouldBe.Equal("USD", res.KeyProduct.Platforms[0].Price.Currency)
	shouldBe.False(res.KeyProduct.Platforms[0].Price.IsFallback)

	res = billingpb.GetKeyProductInfoResponse{}
	err = suite.service.GetKeyProductInfo(context.TODO(), &billingpb.GetKeyProductInfoRequest{Currency: "EUR", KeyProductId: response.Product.Id, Language: "ru"}, &res)
	shouldBe.Nil(err)
	shouldBe.Nil(res.Message)
	shouldBe.NotNil(res.KeyProduct)
	shouldBe.Equal(response.Product.Id, res.KeyProduct.Id)
	shouldBe.Equal(initialName, res.KeyProduct.Name)
	shouldBe.Equal("blah-blah-blah", res.KeyProduct.Description)
	shouldBe.Equal(1, len(res.KeyProduct.Platforms))
	shouldBe.Equal("steam", res.KeyProduct.Platforms[0].Id)
	shouldBe.EqualValues(20, res.KeyProduct.Platforms[0].Price.Amount)
	shouldBe.Equal("EUR", res.KeyProduct.Platforms[0].Price.Currency)
	shouldBe.False(res.KeyProduct.Platforms[0].Price.IsFallback)

	res = billingpb.GetKeyProductInfoResponse{}
	err = suite.service.GetKeyProductInfo(context.TODO(), &billingpb.GetKeyProductInfoRequest{Currency: "UNK", KeyProductId: response.Product.Id, Language: "ru"}, &res)
	shouldBe.Nil(err)
	shouldBe.Nil(res.Message)
	shouldBe.NotNil(res.KeyProduct)
	shouldBe.Equal(response.Product.Id, res.KeyProduct.Id)
	shouldBe.Equal(initialName, res.KeyProduct.Name)
	shouldBe.Equal("blah-blah-blah", res.KeyProduct.Description)
	shouldBe.Equal(1, len(res.KeyProduct.Platforms))
	shouldBe.Equal("steam", res.KeyProduct.Platforms[0].Id)
	shouldBe.EqualValues(10, res.KeyProduct.Platforms[0].Price.Amount)
	shouldBe.Equal("USD", res.KeyProduct.Platforms[0].Price.Currency)
	shouldBe.True(res.KeyProduct.Platforms[0].Price.IsFallback)

	res = billingpb.GetKeyProductInfoResponse{}
	err = suite.service.GetKeyProductInfo(context.TODO(), &billingpb.GetKeyProductInfoRequest{Currency: "RUB", KeyProductId: response.Product.Id, Language: "ru"}, &res)
	shouldBe.Nil(err)
	shouldBe.Nil(res.Message)
	shouldBe.NotNil(res.KeyProduct)
	shouldBe.Equal(response.Product.Id, res.KeyProduct.Id)
	shouldBe.Equal(initialName, res.KeyProduct.Name)
	shouldBe.Equal("blah-blah-blah", res.KeyProduct.Description)
	shouldBe.Equal(1, len(res.KeyProduct.Platforms))
	shouldBe.Equal("steam", res.KeyProduct.Platforms[0].Id)
	shouldBe.EqualValues(10, res.KeyProduct.Platforms[0].Price.Amount)
	shouldBe.Equal("USD", res.KeyProduct.Platforms[0].Price.Currency)
	shouldBe.True(res.KeyProduct.Platforms[0].Price.IsFallback)

	res = billingpb.GetKeyProductInfoResponse{}
	err = suite.service.GetKeyProductInfo(context.TODO(), &billingpb.GetKeyProductInfoRequest{Country: "RUS", KeyProductId: response.Product.Id, Language: "ru"}, &res)
	shouldBe.Nil(err)
	shouldBe.Nil(res.Message)
	shouldBe.NotNil(res.KeyProduct)
	shouldBe.Equal(response.Product.Id, res.KeyProduct.Id)
	shouldBe.Equal(initialName, res.KeyProduct.Name)
	shouldBe.Equal("blah-blah-blah", res.KeyProduct.Description)
	shouldBe.Equal(1, len(res.KeyProduct.Platforms))
	shouldBe.Equal("steam", res.KeyProduct.Platforms[0].Id)
	shouldBe.EqualValues(10, res.KeyProduct.Platforms[0].Price.Amount)
	shouldBe.Equal("USD", res.KeyProduct.Platforms[0].Price.Currency)
	shouldBe.True(res.KeyProduct.Platforms[0].Price.IsFallback)
}

func (suite *KeyProductTestSuite) Test_GetPlatforms() {
	shouldBe := require.New(suite.T())

	rsp := &billingpb.ListPlatformsResponse{}
	shouldBe.Nil(suite.service.GetPlatforms(context.TODO(), &billingpb.ListPlatformsRequest{
		Limit:  100,
		Offset: 0,
	}, rsp))
	shouldBe.EqualValues(200, rsp.Status)
	shouldBe.NotEmpty(rsp.Platforms)
	shouldBe.EqualValues(9, len(rsp.Platforms))
	shouldBe.Equal(rsp.Platforms[0], availablePlatforms["steam"])

	rsp = &billingpb.ListPlatformsResponse{}
	shouldBe.Nil(suite.service.GetPlatforms(context.TODO(), &billingpb.ListPlatformsRequest{
		Limit:  1,
		Offset: 0,
	}, rsp))
	shouldBe.EqualValues(200, rsp.Status)
	shouldBe.Equal(1, len(rsp.Platforms))

	rsp = &billingpb.ListPlatformsResponse{}
	shouldBe.Nil(suite.service.GetPlatforms(context.TODO(), &billingpb.ListPlatformsRequest{
		Limit:  100,
		Offset: 100,
	}, rsp))
	shouldBe.EqualValues(200, rsp.Status)
	shouldBe.Empty(rsp.Platforms)
}

func (suite *KeyProductTestSuite) Test_GetKeyProduct() {
	shouldBe := require.New(suite.T())

	req := &billingpb.CreateOrUpdateKeyProductRequest{
		Object:          "product",
		Sku:             "ru_double_yeti",
		Name:            map[string]string{"en": initialName},
		DefaultCurrency: "USD",
		Description:     map[string]string{"en": "blah-blah-blah"},
		LongDescription: map[string]string{"en": "Super game steam keys"},
		Url:             "http://test.ru/dffdsfsfs",
		Cover: &billingpb.ImageCollection{
			UseOneForAll: false,
			Images: &billingpb.LocalizedUrl{
				En: "/home/image.jpg",
			},
		},
		MerchantId: merchantId,
		ProjectId:  projectId,
		Metadata: map[string]string{
			"SomeKey": "SomeValue",
		},
	}

	response := billingpb.KeyProductResponse{}
	err := suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, &response)
	shouldBe.Nil(err)
	shouldBe.Nil(response.Message)
	res := response.Product

	response = billingpb.KeyProductResponse{}
	err = suite.service.GetKeyProduct(context.TODO(), &billingpb.RequestKeyProductMerchant{Id: res.Id, MerchantId: res.MerchantId}, &response)
	shouldBe.Nil(err)
	shouldBe.Nil(response.Message)

	product := response.Product

	shouldBe.Equal(res.Name["en"], product.Name["en"])
	shouldBe.Equal(res.DefaultCurrency, product.DefaultCurrency)
	shouldBe.Equal(res.Sku, product.Sku)
	shouldBe.Equal(res.Object, product.Object)
	shouldBe.Equal(res.Enabled, product.Enabled)
	shouldBe.Equal(res.Description, product.Description)
	shouldBe.Equal(res.LongDescription, product.LongDescription)
	shouldBe.Equal(res.Url, product.Url)
	shouldBe.Equal(res.Cover.Images.En, product.Cover.Images.En)
	shouldBe.Equal(res.Metadata, product.Metadata)
	shouldBe.NotNil(product.UpdatedAt)
	shouldBe.NotNil(product.CreatedAt)
	shouldBe.Nil(product.PublishedAt)
	shouldBe.False(product.Enabled)

	err = suite.service.GetKeyProduct(context.TODO(), &billingpb.RequestKeyProductMerchant{Id: res.Id, MerchantId: res.MerchantId}, &response)
	shouldBe.Nil(err)
	shouldBe.Nil(response.Message)

	err = suite.service.GetKeyProduct(context.TODO(), &billingpb.RequestKeyProductMerchant{Id: res.Id, MerchantId: res.MerchantId}, &response)
	shouldBe.Nil(err)
	shouldBe.Nil(response.Message)

	err = suite.service.GetKeyProduct(context.TODO(), &billingpb.RequestKeyProductMerchant{Id: res.Id, MerchantId: res.MerchantId}, &response)
	shouldBe.Nil(err)
	shouldBe.Nil(response.Message)
}

func (suite *KeyProductTestSuite) Test_CreateOrUpdateKeyProduct() {
	shouldBe := require.New(suite.T())

	req := &billingpb.CreateOrUpdateKeyProductRequest{
		Object:          "product",
		Sku:             "ru_double_yeti",
		Name:            map[string]string{"en": initialName},
		DefaultCurrency: "USD",
		Description:     map[string]string{"en": "blah-blah-blah"},
		LongDescription: map[string]string{"en": "Super game steam keys"},
		Url:             "http://test.ru/dffdsfsfs",
		Cover: &billingpb.ImageCollection{
			UseOneForAll: false,
			Images: &billingpb.LocalizedUrl{
				En: "/home/image.jpg",
			},
		},
		MerchantId: merchantId,
		ProjectId:  projectId,
		Metadata: map[string]string{
			"SomeKey": "SomeValue",
		},
	}

	response := billingpb.KeyProductResponse{}
	err := suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, &response)
	res := response.Product

	shouldBe.Nil(err)
	shouldBe.EqualValuesf(200, response.Status, "%s", response.Message)
	shouldBe.Equal(res.Name["en"], req.Name["en"])
	shouldBe.Equal(res.DefaultCurrency, req.DefaultCurrency)
	shouldBe.Equal(res.Sku, req.Sku)
	shouldBe.Equal(res.Object, req.Object)
	shouldBe.Equal(res.Description, req.Description)
	shouldBe.Equal(res.LongDescription, req.LongDescription)
	shouldBe.Equal(res.Url, req.Url)
	shouldBe.Equal(res.Cover.Images.En, req.Cover.Images.En)
	shouldBe.Equal(res.Metadata, req.Metadata)
	shouldBe.NotNil(res.UpdatedAt)
	shouldBe.NotNil(res.CreatedAt)
	shouldBe.Nil(res.PublishedAt)
	shouldBe.False(res.Enabled)
	shouldBe.NotEmpty(res.Id)

	req.Id = res.Id
	res2 := billingpb.KeyProductResponse{}
	err = suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, &res2)
	shouldBe.Nil(err)
	shouldBe.Nil(res2.Message)

	res2 = billingpb.KeyProductResponse{}
	req.Id = primitive.NewObjectID().Hex()
	err = suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, &res2)
	shouldBe.Nil(err)
	shouldBe.NotNil(res2.Message)
	shouldBe.EqualValuesf(400, res2.Status, "%s", res2.Message)

	req.Sku = "NEW SKU"
	req.Id = res.Id
	res2 = billingpb.KeyProductResponse{}
	err = suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, &res2)
	shouldBe.Nil(err)
	shouldBe.NotNil(res2.Message)
	shouldBe.EqualValues(400, res2.Status)

	req.Sku = res.Sku
	req.MerchantId = primitive.NewObjectID().Hex()
	res2 = billingpb.KeyProductResponse{}
	err = suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, &res2)
	shouldBe.Nil(err)
	shouldBe.NotNil(res2.Message)
	shouldBe.EqualValues(400, res2.Status)

	req.MerchantId = res.MerchantId
	req.ProjectId = primitive.NewObjectID().Hex()
	res2 = billingpb.KeyProductResponse{}
	err = suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, &res2)
	shouldBe.Nil(err)
	shouldBe.NotNil(res2.Message)
	shouldBe.EqualValues(billingpb.ResponseStatusSystemError, res2.Status)
	shouldBe.Equal(keyProductInternalError, res2.Message)
}

func (suite *KeyProductTestSuite) Test_GetKeyProducts() {
	shouldBe := require.New(suite.T())

	req := &billingpb.ListKeyProductsRequest{
		MerchantId: merchantId,
		ProjectId:  projectId,
	}
	res := &billingpb.ListKeyProductsResponse{}
	err := suite.service.GetKeyProducts(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(0, res.Count)
	shouldBe.EqualValues(0, res.Offset)
	shouldBe.EqualValues(0, len(res.Products))

	for i := 0; i < 10; i++ {
		suite.createKeyProduct()
	}

	err = suite.service.GetKeyProducts(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.EqualValuesf(200, res.Status, "%s", res.Message)
	shouldBe.EqualValues(10, res.Count)
	shouldBe.EqualValues(0, res.Offset)
	shouldBe.EqualValues(10, len(res.Products))
	shouldBe.EqualValues(0, res.Products[0].Platforms[0].Count)

	req.Offset = 9
	err = suite.service.GetKeyProducts(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(10, res.Count)
	shouldBe.EqualValues(1, len(res.Products))

	req.Offset = 0
	req.Limit = 2
	err = suite.service.GetKeyProducts(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(10, res.Count)
	shouldBe.EqualValues(2, len(res.Products))

	req.Offset = 0
	req.Limit = 0
	req.Sku = "some sku"
	req.Name = "some name"
	err = suite.service.GetKeyProducts(context.TODO(), req, res)
	shouldBe.Nil(err)

	req.Offset = 0
	req.Limit = 100
	req.Sku = ""
	req.Name = ""
	req.Enabled = "true"
	err = suite.service.GetKeyProducts(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(0, len(res.Products))

	req.Offset = 0
	req.Limit = 100
	req.Sku = ""
	req.Name = ""
	req.Enabled = "false"
	err = suite.service.GetKeyProducts(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(10, len(res.Products))
}

func (suite *KeyProductTestSuite) getKeyProduct(id string) *billingpb.KeyProduct {
	suite.T().Helper()

	res := &billingpb.KeyProductResponse{}
	err := suite.service.GetKeyProduct(context.TODO(), &billingpb.RequestKeyProductMerchant{MerchantId: merchantId, Id: id}, res)
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), res.Message)
	return res.Product
}

func (suite *KeyProductTestSuite) createKeyProduct() *billingpb.KeyProduct {
	suite.T().Helper()

	req := &billingpb.CreateOrUpdateKeyProductRequest{
		Object:          "product",
		Sku:             primitive.NewObjectID().Hex(),
		Name:            map[string]string{"en": initialName},
		DefaultCurrency: "USD",
		Description:     map[string]string{"en": "blah-blah-blah"},
		LongDescription: map[string]string{"en": "Super game steam keys"},
		Url:             "http://test.ru/dffdsfsfs",
		Cover: &billingpb.ImageCollection{
			UseOneForAll: false,
			Images: &billingpb.LocalizedUrl{
				En: "/home/image.jpg",
			},
		},
		MerchantId: merchantId,
		ProjectId:  projectId,
		Platforms: []*billingpb.PlatformPrice{
			{
				Id: "steam",
				Prices: []*billingpb.ProductPrice{
					{Region: "USD", Currency: "USD", Amount: 66.66},
				},
			},
		},
		Metadata: map[string]string{
			"SomeKey": "SomeValue",
		},
	}

	res := &billingpb.KeyProductResponse{}
	err := suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, res)
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), res.Message)
	return res.Product
}

func (suite *KeyProductTestSuite) Test_UpdatePlatformPrices_WithBadPrice_Error() {
	shouldBe := require.New(suite.T())
	product := suite.createKeyProduct()
	req := &billingpb.CreateOrUpdateKeyProductRequest{
		Id:              product.Id,
		MerchantId:      product.MerchantId,
		Name:            product.Name,
		Description:     product.Description,
		ProjectId:       product.MerchantId,
		DefaultCurrency: product.DefaultCurrency,
		Platforms: []*billingpb.PlatformPrice{
			{
				Id: "steam",
				Prices: []*billingpb.ProductPrice{
					{Region: "USD", Currency: "RUB", Amount: 66.66},
				},
			},
		},
	}

	res := &billingpb.KeyProductResponse{}
	err := suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(400, res.Status)
	shouldBe.NotEmpty(res.Message)
}

func (suite *KeyProductTestSuite) Test_UpdatePlatformPrices() {
	shouldBe := require.New(suite.T())
	product := suite.createKeyProduct()
	req := &billingpb.CreateOrUpdateKeyProductRequest{
		Id:              product.Id,
		MerchantId:      product.MerchantId,
		Name:            product.Name,
		Description:     product.Description,
		ProjectId:       product.MerchantId,
		DefaultCurrency: product.DefaultCurrency,
		Platforms: []*billingpb.PlatformPrice{
			{
				Id: "steam",
				Prices: []*billingpb.ProductPrice{
					{Region: "USD", Currency: "USD", Amount: 66.66},
				},
			},
		},
	}

	res := &billingpb.KeyProductResponse{}
	err := suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.Nil(res.Message)

	prices := res.Product.Platforms[0].Prices
	shouldBe.Equal(1, len(prices))
	shouldBe.Equal(66.66, prices[0].Amount)
	shouldBe.Equal("USD", prices[0].Currency)

	req = &billingpb.CreateOrUpdateKeyProductRequest{
		Id:              product.Id,
		MerchantId:      product.MerchantId,
		ProjectId:       product.MerchantId,
		DefaultCurrency: product.DefaultCurrency,
		Name:            product.Name,
		Description:     product.Description,
		Platforms: []*billingpb.PlatformPrice{
			{
				Id: "steam",
				Prices: []*billingpb.ProductPrice{
					{Region: "USD", Currency: "USD", Amount: 77.66},
					{Region: "EUR", Currency: "EUR", Amount: 77.77},
				},
			},
		},
	}

	res = &billingpb.KeyProductResponse{}
	err = suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, res)

	shouldBe.Nil(err)
	shouldBe.Nil(res.Message)

	req = &billingpb.CreateOrUpdateKeyProductRequest{
		Id:              product.Id,
		MerchantId:      product.MerchantId,
		ProjectId:       product.MerchantId,
		DefaultCurrency: product.DefaultCurrency,
		Name:            product.Name,
		Description:     product.Description,
		Platforms: []*billingpb.PlatformPrice{
			{
				Id:            "best_store_ever",
				EulaUrl:       "http://www.example.com",
				ActivationUrl: "http://www.example.com",
				Prices: []*billingpb.ProductPrice{
					{Region: "RUB", Currency: "RUB", Amount: 0.01},
					{Region: "USD", Currency: "USD", Amount: 66.66},
				},
			},
		},
	}

	res = &billingpb.KeyProductResponse{}
	err = suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.NotNil(res.Message)
	shouldBe.EqualValues(400, res.Status)

	req = &billingpb.CreateOrUpdateKeyProductRequest{
		Id:              product.Id,
		MerchantId:      product.MerchantId,
		ProjectId:       product.MerchantId,
		DefaultCurrency: product.DefaultCurrency,
		Name:            product.Name,
		Description:     product.Description,
		Platforms: []*billingpb.PlatformPrice{
			{
				Id:            "best_store_ever",
				EulaUrl:       "http://www.example.com",
				ActivationUrl: "http://www.example.com",
				Prices: []*billingpb.ProductPrice{
					{Region: "RUB", Currency: "RUB", Amount: 0.01},
					{Region: "USD", Currency: "USD", Amount: 66.66},
				},
			},
			{
				Id:            "another_best_store_ever",
				EulaUrl:       "http://www.example.com",
				ActivationUrl: "http://www.example.com",
				Prices: []*billingpb.ProductPrice{
					{Region: "RUB", Currency: "RUB", Amount: 0.01},
					{Region: "USD", Currency: "USD", Amount: 66.66},
				},
			},
		},
	}
	res = &billingpb.KeyProductResponse{}
	err = suite.service.CreateOrUpdateKeyProduct(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.NotNil(res.Message)
	shouldBe.EqualValues(400, res.Status)
}

func (suite *KeyProductTestSuite) Test_PublishKeyProduct() {
	shouldBe := require.New(suite.T())

	product := suite.createKeyProduct()
	req := &billingpb.PublishKeyProductRequest{
		KeyProductId: product.Id,
		MerchantId:   merchantId,
	}
	res := &billingpb.KeyProductResponse{}
	err := suite.service.PublishKeyProduct(context.TODO(), req, res)
	shouldBe.Nil(err)
	shouldBe.Nil(res.Message)
	shouldBe.EqualValues(200, res.Status)
	shouldBe.True(res.Product.Enabled)
	shouldBe.NotNil(res.Product.PublishedAt)
}

func (suite *KeyProductTestSuite) Test_DeleteKeyProduct() {
	shouldBe := require.New(suite.T())
	product := suite.createKeyProduct()

	res := &billingpb.EmptyResponseWithStatus{}
	err := suite.service.DeleteKeyProduct(context.TODO(), &billingpb.RequestKeyProductMerchant{Id: product.Id, MerchantId: merchantId}, res)
	shouldBe.Nil(err)
	shouldBe.Nil(res.Message)

	res = &billingpb.EmptyResponseWithStatus{}
	err = suite.service.DeleteKeyProduct(context.TODO(), &billingpb.RequestKeyProductMerchant{Id: product.Id, MerchantId: merchantId}, res)
	shouldBe.Nil(err)
	shouldBe.NotNil(res.Message)
}

func (suite *KeyProductTestSuite) Test_UploadKey() {
	shouldBe := require.New(suite.T())
	product := suite.createKeyProduct()
	fileContent := fmt.Sprintf("%s-%s-%s-%s", RandomString(4), RandomString(4), RandomString(4), RandomString(4))
	file := []byte(fileContent)

	keysRsp := &billingpb.PlatformKeysFileResponse{}
	keysReq := &billingpb.PlatformKeysFileRequest{
		KeyProductId: product.Id,
		PlatformId:   "steam",
		MerchantId:   product.MerchantId,
		File:         file,
	}
	shouldBe.NoError(suite.service.UploadKeysFile(context.TODO(), keysReq, keysRsp))
	shouldBe.EqualValues(200, keysRsp.Status)
	shouldBe.EqualValues(1, keysRsp.TotalCount)
	shouldBe.EqualValues(1, keysRsp.KeysProcessed)

	keysRsp = &billingpb.PlatformKeysFileResponse{}
	keysReq = &billingpb.PlatformKeysFileRequest{
		KeyProductId: product.Id,
		PlatformId:   "steam",
		MerchantId:   product.MerchantId,
		File:         file,
	}
	shouldBe.NoError(suite.service.UploadKeysFile(context.TODO(), keysReq, keysRsp))
	shouldBe.EqualValues(200, keysRsp.Status)
	shouldBe.EqualValues(1, keysRsp.TotalCount)
	shouldBe.EqualValues(0, keysRsp.KeysProcessed)
}

func (suite *KeyProductTestSuite) Test_CheckSkuAndKeyProject() {
	shouldBe := require.New(suite.T())
	rsp := &billingpb.EmptyResponseWithStatus{}
	err := suite.service.CheckSkuAndKeyProject(context.TODO(), &billingpb.CheckSkuAndKeyProjectRequest{ProjectId: projectId, Sku: "TEST_SKU"}, rsp)
	shouldBe.NoError(err)
	shouldBe.EqualValues(200, rsp.Status)

	product := suite.createKeyProduct()
	err = suite.service.CheckSkuAndKeyProject(context.TODO(), &billingpb.CheckSkuAndKeyProjectRequest{ProjectId: product.ProjectId, Sku: product.Sku}, rsp)
	shouldBe.NoError(err)
	shouldBe.EqualValues(400, rsp.Status)
	shouldBe.NotNil(rsp.Message)
}

func (suite *KeyProductTestSuite) Test_UnPublishKeyProduct() {
	shouldBe := require.New(suite.T())
	product := suite.createKeyProduct()

	res := &billingpb.KeyProductResponse{}
	err := suite.service.PublishKeyProduct(context.TODO(), &billingpb.PublishKeyProductRequest{KeyProductId: product.Id, MerchantId: product.MerchantId}, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(200, res.Status)

	err = suite.service.UnPublishKeyProduct(context.TODO(), &billingpb.UnPublishKeyProductRequest{KeyProductId: product.Id}, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(400, res.Status)

	err = suite.service.UnPublishKeyProduct(context.TODO(), &billingpb.UnPublishKeyProductRequest{KeyProductId: product.Id, MerchantId: product.MerchantId}, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(200, res.Status)

	err = suite.service.UnPublishKeyProduct(context.TODO(), &billingpb.UnPublishKeyProductRequest{KeyProductId: product.Id, MerchantId: product.MerchantId}, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(400, res.Status)
	shouldBe.Equal(keyProductNotPublished, res.Message)

	emptyRes := &billingpb.EmptyResponseWithStatus{}
	err = suite.service.DeleteKeyProduct(context.TODO(), &billingpb.RequestKeyProductMerchant{Id: product.Id, MerchantId: product.MerchantId}, emptyRes)
	shouldBe.Nil(err)
	shouldBe.EqualValuesf(200, emptyRes.Status, "%v", emptyRes.Message)

	err = suite.service.UnPublishKeyProduct(context.TODO(), &billingpb.UnPublishKeyProductRequest{KeyProductId: product.Id, MerchantId: product.MerchantId}, res)
	shouldBe.Nil(err)
	shouldBe.EqualValues(400, res.Status)
	shouldBe.Equal(keyProductNotFound, res.Message)
}
