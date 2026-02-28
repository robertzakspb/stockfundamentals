package main

import (
	"context"
	"net/http"
	"time"

	api_security "github.com/compoundinvest/stockfundamentals/internal/interface/api/security"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/interface/api/time-series"
	"github.com/ydb-platform/ydb-go-sdk/v3"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/dataseed"

	"github.com/compoundinvest/stockfundamentals/internal/interface/api/account/portfolio"
	dividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())
	addEndpoints(router)

	router.Run("localhost:8080")
}

func addEndpoints(router *gin.Engine) {
	router.POST("/migration/initial-seed", dataseed.InitialSeed)

	router.GET("/health-check", healthCheck)

	router.GET("/portfolio", portfolio.GetPortfolio)
	router.GET("/account-portfolio", portfolio.GetAccountPortfolio)
	router.POST("/update-portfolio", portfolio.UpdatePortfolio)

	router.POST("/fetch/dividends", dividend.StartDividendFetchingJob)
	router.GET("/all-dividends", dividend.GetAllDividends)
	router.GET("/upcoming-dividends", dividend.GetUpcomingDividends) //TODO: Deprecate (may now be simulated via /all-dividends)

	router.POST("/fetch/securities", api_security.StartSecurityMasterImportJob)

	router.POST("/fetch/time-series", timeseries.StartTimeSeriesImportJob)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {

	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)
	if err != nil {
	}
	dataseed.CreateDividendForecastTable(ctx, db, db.Table())
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "The service is running without any issues")
}
