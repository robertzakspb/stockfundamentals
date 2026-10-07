package cacheorchestrator

import (
	"github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/cache/cache"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/quoteservice"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func LoadQuoteCache() {
	err := loadStockQuoteCache()
	if err != nil {
		logger.LogError(err)
	}
	err = loadBondQuoteCache()
	if err != nil {
		logger.LogError(err)
	}
	logger.Log("Cached stock quotes and bond quotes", logger.INFORMATION)
}

func loadStockQuoteCache() error {
	stocks, err := security_master.GetAllSecuritiesFromDB()
	if err != nil {
		return err
	}
	securities := security_master.ConvertStocksToSecurities(stocks)

	quotes, err := quoteservice.FetchInternationalStockQuotes(securities)
	if err != nil {
		return err
	}

	cache.UpdateStockCache(quotes)

	return nil
}

func loadBondQuoteCache() error {
	bondList, err := bondservice.GetAllBonds()
	if err != nil {
		return err
	}

	figis := make([]string, len(bondList))
	for i := range bondList {
		figis = append(figis, bondList[i].Figi)
	}
	quotes, err := quoteservice.FetchBondQuotesFromTapi(figis)

	cache.UpdateBondsCache(quotes)

	return nil
}
