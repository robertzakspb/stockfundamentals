package jobs

import (
	"net/http"
	"sync"

	// accountreturnapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/account-return"
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
	wg := sync.WaitGroup{}
	wg.Go(func() { api_security.StartSecurityMasterImportJob(c) })
	wg.Wait() //Need to import the security master before updating the portfolio

	go portfolio.UpdatePortfolio(c)

	wg.Go(func() { forexapi.StartForexImportJob(c) })
	wg.Wait() //Need to fetch the latest forex rates before proceeding to update the bonds' ACI

	bondsapi.StartBondAccruedInterestUpdateJob(c)

	//accountreturnapi.StartMarketValueSnapshotJob(c) //FIXME

	c.JSON(http.StatusOK, "Daily jobs were successfully started/executed")
}

func StartHeavyJobs(c *gin.Context) {
	timeseries.StartTimeSeriesImportJob(c)  //Completes in 18 minutes if run separately
	apidividend.StartDividendFetchingJob(c) //Completes in 1.5 minutes if run separately
	bondsapi.StartBondAndCouponImportJob(c) //Completes in 18.5 minutes if run separately
	c.JSON(http.StatusOK, "Heavy jobs were successfully started/executed")
}

