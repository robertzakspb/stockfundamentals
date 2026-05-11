package jobs

import (
	"net/http"
	"sync"

	// accountreturnapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/account-return"
	accountreturnapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/account-return"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/stock-portfolio"
	bondsapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/bonds"
	forexapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/forex"
	apidividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	api_security "github.com/compoundinvest/stockfundamentals/internal/interface/api/security"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/interface/api/time-series"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func StartAllJobs(c *gin.Context) {
	StartDailyJobs(c)
	StartHeavyJobs(c)
	c.JSON(http.StatusOK, shared.StringResponse{Message: "All jobs were successfully started/executed"})
}

func StartDailyJobs(c *gin.Context) {
	wg := sync.WaitGroup{}

	wg.Go(func() { api_security.StartSecurityMasterImportJob(c) })
	wg.Go(func() { forexapi.StartForexImportJob(c) })

	wg.Wait() //Need to fetch the latest forex rates before proceeding to update the bonds' ACI

	bondsapi.StartBondAccruedInterestUpdateJob(c)

	portfolio.UpdatePortfolio(c)

	accountreturnapi.StartMarketValueSnapshotJob(c)

	c.JSON(http.StatusOK, shared.StringResponse{Message: "Daily jobs were successfully started/executed"})
}

func StartHeavyJobs(c *gin.Context) {
	timeseries.StartTimeSeriesImportJob(c)  //Completes in 18 minutes if run separately
	apidividend.StartDividendFetchingJob(c) //Completes in 1.5 minutes if run separately
	bondsapi.StartBondAndCouponImportJob(c) //Completes in 18.5 minutes if run separately
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Heavy jobs were successfully started/executed"})
}
