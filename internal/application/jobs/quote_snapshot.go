package jobservice

import (
	"sync"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/quoteservice"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	timeseriesdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func ExecuteQuoteSnapshotJob() error {

	stocks, err := security_master.GetAllSecuritiesFromDB()
	if err != nil {
		return err
	}
	stockFigis := security_master.ExtractFigisFromSecurities(stocks)

	bondList, err := bondservice.GetAllBonds()
	if err != nil {
		return err
	}
	bondFigis := bondservice.ExtractBondFigis(&bondList)

	stockQuotes := []entity.SimpleQuote{}
	bondQuotes := []entity.BondQuote{}

	wg := sync.WaitGroup{}
	wg.Go(func() {
		bondQuotes, err = quoteservice.FetchBondQuotes(bondFigis)
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
		}
	})
	wg.Go(func() {
		stockQuotes, err = quoteservice.FetchStockQuotes(stockFigis)
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
		}
	})
	wg.Wait()

	if err != nil {
		return err
	}

	err = timeseriesdb.SaveTimeSeriesToDB(&stockQuotes)
	if err != nil {
		return err
	}

	bondDbQuotes := quoteservice.MapBondQuotesToDbModels(bondQuotes)
	err = timeseriesdb.SaveBondQuotes(bondDbQuotes)
	if err != nil {
		return err
	}

	return nil
}
