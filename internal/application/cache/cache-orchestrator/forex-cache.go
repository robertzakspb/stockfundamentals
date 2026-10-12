package cacheorchestrator

import (
	"errors"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/application/cache/cache"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func LoadForexCache() error {
	if cache.ForexCacheIsLoaded() {
		err := errors.New("Cache is already populated")
		logger.LogError(err)
		return err
	}

	supportedCurrencyPairs := forexservice.SupportedCurrencyPairs()
	rates, err := forexservice.GetExchangeRates(supportedCurrencyPairs, time.Now())
	if err != nil {
		logger.LogError(err)
	}
	//Because the DB could contain several (but not all) exchange rates for today, we should ensure that we fetch the remaining ones before populating the cache
	if len(rates) != len(supportedCurrencyPairs) {
		forexservice.ImportForexRates()
		rates, err = forexservice.GetExchangeRates(supportedCurrencyPairs, time.Now())
	}

	cache.UpdateForexCache(rates)

	logger.Log("Successfuly loaded the forex rate cache", logger.INFORMATION)
	return nil
}
