package main

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/dataseed"
	"github.com/compoundinvest/stockfundamentals/internal/interface/api/jobs"
	api_security "github.com/compoundinvest/stockfundamentals/internal/interface/api/security"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/interface/api/time-series"

	bondsapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/bond-portfolio"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/stock-portfolio"
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
	router.POST("dividend/forecast", dividend.CreateNewDividendForecast)
	router.GET("/dividend/forecasts", dividend.GetDividendForecasts)
	router.GET("dividend/forecasts-grouped-by-security", dividend.GetDividendForecastsGroupedBySecurity)

	router.POST("jobs/import-bonds-and-coupons", jobs.StartBondAndCouponImportJob)
	router.POST("bonds/new-position-lot", bondsapi.AddBondPositionLotToAccount)
	router.GET("bonds/position-lots", bondsapi.GetAccountPositionLots)

	router.POST("/fetch/securities", api_security.ExecuteSecurityMasterImportJob)

	router.POST("/fetch/time-series", timeseries.StartTimeSeriesImportJob)

	router.POST("/start-all-jobs", jobs.StartAllJobs)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "The service is running without any issues")
}
