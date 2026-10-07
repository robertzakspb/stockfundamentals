package cache

import (
	"errors"
	"sync"

	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

var forexRates []forexservice.ForexRate
var fxMutex sync.Mutex

func UpdateForexCache(rates []forexservice.ForexRate) {
	fxMutex.Lock()
	defer fxMutex.Unlock()
	forexRates = rates
}

func GetForexRateCacheForCurrencyPairs(pairs ...string) ([]forexservice.ForexRate, error) {
	fxMutex.Lock()
	defer fxMutex.Unlock()

	if len(forexRates) == 0 {
		err := errors.New("Attempting to fetch the forex cache despite its being empty")
		logger.LogError(err)
		return []forexservice.ForexRate{}, err
	}

	targetRates, err := forexservice.СollapseRatesIntoTargetCrossRates(pairs, forexRates)
	if err != nil {
		logger.LogError(err)
	}
	return targetRates, err
}

func GetForexRateCacheForCurrencyPair(pair string) (forexservice.ForexRate, error) {
	rates, err := GetForexRateCacheForCurrencyPairs(pair)
	if err != nil {
		return forexservice.ForexRate{}, err
	}
	if len(rates) == 0 {
		return forexservice.ForexRate{}, errors.New("Failed to find a cached forex rate for " + pair)
	}
	return rates[0], nil
}

func ForexCacheIsLoaded() bool {
	fxMutex.Lock()
	defer fxMutex.Unlock()

	return len(forexRates) > 0
}
