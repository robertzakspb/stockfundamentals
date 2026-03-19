package jobs

import (
	"net/http"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/stock-portfolio"
	apidividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	api_security "github.com/compoundinvest/stockfundamentals/internal/interface/api/security"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/interface/api/time-series"
	"github.com/gin-gonic/gin"
)

func StartBondAndCouponImportJob(c *gin.Context) {
	go bondservice.ImportAllBondsAndCoupons()

	c.JSON(http.StatusOK, "The bond import job has been successfully started")
}

func StartAllJobs(c *gin.Context) {
	//Ensure that this job is synchronously started and completed before any other jobs run to prevent collisions
	api_security.ExecuteSecurityMasterImportJob(c)

	portfolio.UpdatePortfolio(c)

	timeseries.StartTimeSeriesImportJob(c)
	apidividend.StartDividendFetchingJob(c)
	StartBondAndCouponImportJob(c)

	c.JSON(http.StatusOK, "All jobs were successfully started/executed")
}
