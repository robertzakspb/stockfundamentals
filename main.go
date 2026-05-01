package main

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/dataseed"
	bondsapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/bonds"
	forexapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/forex"
	"github.com/compoundinvest/stockfundamentals/internal/interface/api/jobs"
	api_security "github.com/compoundinvest/stockfundamentals/internal/interface/api/security"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/interface/api/time-series"

	accountreturnapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/account-return"
	bondportfolioapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/bond-portfolio"
	divcalapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/dividend-calendar"
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
	router.POST("migration/initial-seed", dataseed.InitialSeed)

	router.GET("health-check", healthCheck)

	router.GET("account-portfolio", portfolio.GetAccountPortfolio)
	router.POST("update-portfolio", portfolio.UpdatePortfolio)
	router.GET("account/return", accountreturnapi.GetAccountReturn)
	router.POST("account/save-market-value-snapshots", accountreturnapi.StartMarketValueSnapshotJob)
	router.GET("account/bond-portfolio-analysis", accountreturnapi.GetPortfolioOverview)

	router.POST("fetch/dividends", dividend.StartDividendFetchingJob)
	router.GET("all-dividends", dividend.GetAllDividends)
	router.GET("upcoming-dividends", dividend.GetUpcomingDividends) //TODO: Deprecate (may now be simulated via /all-dividends)
	router.POST("dividend/forecast", dividend.CreateNewDividendForecast)
	router.GET("dividend/forecasts", dividend.GetDividendForecasts)
	router.GET("dividend/forecasts-grouped-by-security", dividend.GetDividendForecastsGroupedBySecurity)

	router.GET("dividend/calendar", divcalapi.GetAccountDividendCalendar)

	router.POST("jobs/import-bonds-and-coupons", bondsapi.StartBondAndCouponImportJob)

	router.GET("bonds/russian-government-bonds", bondsapi.GetRussianGovernmentBondsWithFixedOrConstantCoupon)
	router.GET("bonds/quasi-foreign-bonds", bondsapi.GetQuasiForeignBonds)
	router.POST("bonds/new-position-lot", bondportfolioapi.AddBondPositionLotToAccount)
	router.GET("bonds/position-lots", bondportfolioapi.GetAccountPositionLots)
	router.GET("bonds/account/timeline", bondportfolioapi.GetAccountBondTimeline)
	router.POST("bonds/update-all-bonds-aci", bondsapi.StartBondAccruedInterestUpdateJob)

	router.POST("fetch/securities", api_security.StartSecurityMasterImportJob)

	router.POST("fetch/time-series", timeseries.StartTimeSeriesImportJob)

	router.POST("fetch/fx-rates", forexapi.StartForexImportJob)

	router.POST("start-all-jobs", jobs.StartAllJobs)
	router.POST("start-daily-jobs", jobs.StartDailyJobs)
	router.POST("start-heavy-jobs", jobs.StartHeavyJobs)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "The service is running without any issues")
}
