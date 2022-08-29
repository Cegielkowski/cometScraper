package main

import (
	"cometScraper/tools/scraper/pkg/applicant"
	"cometScraper/tools/scraper/pkg/crawler"
	"net/http"
	"time"

	_ "cometScraper/docs"
	"cometScraper/utils"

	"cometScraper/config"
	httpDelivery "cometScraper/delivery/http"
	appMiddleware "cometScraper/delivery/middleware"
	"cometScraper/infrastructure/datastore"
	pgsqlRepository "cometScraper/repository/pgsql"
	redisRepository "cometScraper/repository/redis"
	"cometScraper/usecase"
	"cometScraper/utils/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Load config
	configApp := config.LoadConfig()

	// Setup logger
	appLogger := logger.NewApiLogger(configApp)
	appLogger.InitLogger()

	// Setup infra
	dbInstance, err := datastore.NewDatabase(configApp.DatabaseURL)
	utils.PanicIfNeeded(err)

	cacheInstance, err := datastore.NewCache(configApp.CacheURL)
	utils.PanicIfNeeded(err)

	// Setup repository
	redisRepo := redisRepository.NewRedisRepository(cacheInstance)
	cometScraperRepo := pgsqlRepository.NewPgsqlCometScraperRepository(dbInstance)

	//Setup Scraper
	applicantInstance := applicant.NewApplicant()
	cometCrawler := crawler.NewCometCrawler(configApp.Elements, applicantInstance)

	// Setup usecase
	cometScraperUC := usecase.NewCometScraperUsecase(cometScraperRepo, redisRepo, cometCrawler)

	// Setup app middleware
	appMiddleware := appMiddleware.NewMiddleware(appLogger)

	// Setup route engine & middleware
	e := echo.New()
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(configApp.ContextTimeout) * time.Second,
	}))
	e.Use(middleware.CORS())
	e.Use(appMiddleware.RequestID())
	e.Use(appMiddleware.Logger())
	e.Use(middleware.Recover())

	// Setup handler
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "i am alive")
	})

	httpDelivery.NewCometScraperHandler(e, cometScraperUC)

	e.Logger.Fatal(e.Start(":" + configApp.ServerPORT))
}
