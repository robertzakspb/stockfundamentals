package main

import (
	"net/http"

	api_security "github.com/compoundinvest/stockfundamentals/internal/interface/api/security"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/interface/api/time-series"

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
	router.GET("/upcoming-dividends", dividend.GetUpcomingDividends)

	router.POST("/fetch/securities", api_security.StartSecurityMasterImportJob)

	router.POST("/fetch/time-series", timeseries.StartTimeSeriesImportJob)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "The service is running without any issues")
}
