package jobservice

import (
	"sync"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	accountmvservice "github.com/compoundinvest/stockfundamentals/internal/application/account/market-value"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/application/timeseries"
)

func StartDailyJobs() {
	wg := sync.WaitGroup{}

	wg.Go(func() { security_master.FetchAndSaveSecurities() })
	wg.Go(func() { forexservice.ImportForexRatesJob() })

	wg.Wait() //Need to fetch the latest forex rates before proceeding to update the bonds' ACI. The security master must also be updated before any other stock-related jobs are executed

	go bondservice.UpdateAllBondsAci()

	portfolio.UpdatePortfolio() //The stock & bonds portfolios must be updated before the market value job so that the latest positions are used in MV calculation
	bondportfolio.ImportTinkoffBondLots()

	go accountmvservice.SaveAccountMarketValueSnapshots()
}

func StartHeavyJobs() {
	go timeseries.FetchAndSaveHistoricalQuotes() //Completes in 18 minutes if run separately
	go appdividend.FetchAndSaveAllDividends()    //Completes in 1.5 minutes if run separately
	go bondservice.ImportAllBondsAndCoupons()    //Completes in 18.5 minutes if run separately
}
