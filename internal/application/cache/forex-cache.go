package cache

import (
	"errors"
	"strconv"
	"sync"

	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

var forexRates []forexservice.ForexRate
var fxMutex sync.Mutex

func LoadForexCache() error {
	fxMutex.Lock()
	defer fxMutex.Unlock()

	if len(forexRates) > 0 {
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

	forexRates = rates
	logger.Log("Successfuly loaded the forex rate cache ("+strconv.Itoa(len(forexRates))+" rates)", logger.INFORMATION)
	return nil
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
