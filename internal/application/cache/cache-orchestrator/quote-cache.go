package cacheorchestrator

import (
	"strconv"
	"sync"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/quoteservice"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

var cachedStockQuotes []entity.SimpleQuote
var stockQuoteMx sync.Mutex

var cachedBondQuotes []entity.BondQuote
var bondQuoteMx sync.Mutex

func LoadQuoteCache() {
	err := loadStockQuoteCache()
	if err != nil {
		logger.LogError(err)
	}
	err = loadBondQuoteCache()
	if err != nil {
		logger.LogError(err)
	}
	logger.Log("Cached "+strconv.Itoa(len(cachedStockQuotes))+" stock quotes and "+strconv.Itoa(len(cachedBondQuotes))+" bond quotes", logger.INFORMATION)
}

func loadStockQuoteCache() error {
	stockQuoteMx.Lock()
	defer stockQuoteMx.Unlock()

	stocks, err := security_master.GetAllSecuritiesFromDB()
	if err != nil {
		return err
	}
	securities := security_master.ConvertStocksToSecurities(stocks)

	quotes, err := quoteservice.FetchInternationalStockQuotes(securities)
	if err != nil {
		return err
	}

	cachedStockQuotes = quotes

	return nil
}

func loadBondQuoteCache() error {
	bondQuoteMx.Lock()
	defer bondQuoteMx.Unlock()

	bondList, err := bondservice.GetAllBonds()
	if err != nil {
		return err
	}

	figis := make([]string, len(bondList))
	for i := range bondList {
		figis = append(figis, bondList[i].Figi)
	}
	quotes, err := quoteservice.FetchBondQuotesFromTapi(figis)

	cachedBondQuotes = quotes

	return nil
}
