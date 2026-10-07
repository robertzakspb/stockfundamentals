package cacheorchestrator

import (
	"errors"

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
	rates, err := forexservice.GetLatestAvailableRates(supportedCurrencyPairs)
	if err != nil {
		logger.LogError(err)
	}
	if len(rates) == 0 {
		forexservice.ImportForexRates()
		rates, err = forexservice.GetLatestAvailableRates(supportedCurrencyPairs)
	}

	cache.UpdateForexCache(rates)

	logger.Log("Successfuly loaded the forex rate cache", logger.INFORMATION)
	return nil
}
