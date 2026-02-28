package jobs

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/interface/api/account/portfolio"
	apidividend "github.com/compoundinvest/stockfundamentals/internal/interface/api/fundamentals/dividend"
	api_security "github.com/compoundinvest/stockfundamentals/internal/interface/api/security"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/interface/api/time-series"
	"github.com/gin-gonic/gin"
)

func StartAllJobs(c *gin.Context) {
	//Ensure that this job is synchronously started and completed before any other jobs run to prevent collisions
	api_security.StartSecurityMasterImportJob(c)

	timeseries.StartTimeSeriesImportJob(c)
	apidividend.StartDividendFetchingJob(c)
	portfolio.UpdatePortfolio(c)

	c.JSON(http.StatusOK, "All jobs were successfully started/executed")
}
