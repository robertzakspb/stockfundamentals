package jobs

import (
	"net/http"
	"sync"

	// accountreturnapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/account-return"
	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	accountmvservice "github.com/compoundinvest/stockfundamentals/internal/application/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/application/timeseries"
	"github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
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

	wg.Go(func() { security_master.FetchAndSaveSecurities() })
	wg.Go(func() { forexservice.ImportForexRatesJob() })

	wg.Wait() //Need to fetch the latest forex rates before proceeding to update the bonds' ACI. The security master must also be updated before any other stock-related jobs are executed

	go bondservice.UpdateAllBondsAci()

	portfolio.UpdatePortfolio() //The stock & bonds portfolios must be updated before the market value job so that the latest positions are used in MV calculation
	bondportfolio.ImportTinkoffBondLots()

	go accountmvservice.SaveAccountMarketValueSnapshots()

	c.JSON(http.StatusOK, shared.StringResponse{Message: "Daily jobs were successfully started/executed"})
}

func StartHeavyJobs(c *gin.Context) {
	go timeseries.FetchAndSaveHistoricalQuotes() //Completes in 18 minutes if run separately
	go appdividend.FetchAndSaveAllDividends()    //Completes in 1.5 minutes if run separately
	go bondservice.ImportAllBondsAndCoupons()    //Completes in 18.5 minutes if run separately
	c.JSON(http.StatusOK, shared.StringResponse{Message: "Heavy jobs were successfully started/executed"})
}
