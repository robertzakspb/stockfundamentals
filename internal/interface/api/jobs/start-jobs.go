package jobs

import (
	"net/http"

	portfolio "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/stock-portfolio"
	bondsapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/bonds"
	forexapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/forex"
	apidividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	api_security "github.com/compoundinvest/stockfundamentals/internal/interface/api/security"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/interface/api/time-series"
	"github.com/gin-gonic/gin"
)

func StartAllJobs(c *gin.Context) {
	//Ensure that this job is synchronously started and completed before any other jobs run to prevent collisions
	StartDailyJobs(c)
	StartHeavyJobs(c)
	c.JSON(http.StatusOK, "All jobs were successfully started/executed")
}

func StartDailyJobs(c *gin.Context) {
	api_security.ExecuteSecurityMasterImportJob(c)
	portfolio.UpdatePortfolio(c)
	bondsapi.UpdateAllBondsAci(c)
	c.JSON(http.StatusOK, "Daily jobs were successfully started/executed")
}

func StartHeavyJobs(c *gin.Context) {
	timeseries.StartTimeSeriesImportJob(c)
	apidividend.StartDividendFetchingJob(c)
	bondsapi.StartBondAndCouponImportJob(c)
	forexapi.StartForexImportJob(c)
	c.JSON(http.StatusOK, "Heavy jobs were successfully started/executed")
}
