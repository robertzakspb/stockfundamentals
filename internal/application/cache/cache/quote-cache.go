package cache

import (
	"errors"
	"sync"

	"github.com/compoundinvest/invest-core/quote/entity"
)

var cachedStockQuotes []entity.SimpleQuote
var stockQuoteMx sync.Mutex

var cachedBondQuotes []entity.BondQuote
var bondQuoteMx sync.Mutex

func UpdateStockAndBondsCache(stockQuotes []entity.SimpleQuote, bondQuotes []entity.BondQuote) {
	stockQuoteMx.Lock()
	bondQuoteMx.Lock()
	defer stockQuoteMx.Unlock()
	defer bondQuoteMx.Unlock()
	if len(stockQuotes) != 0 {
		cachedStockQuotes = stockQuotes
	}
	if len(bondQuotes) != 0 {
		cachedBondQuotes = bondQuotes
	}
}

func GetCachedStockAndBondQuotes(stockFigis, bondFigis []string) ([]entity.SimpleQuote, []entity.BondQuote, error) {
	var stockQuotes []entity.SimpleQuote
	var bondQuotes []entity.BondQuote
	var err error
	var wg sync.WaitGroup

	wg.Go(func() {
		stockQuotes, err = GetCachedStockQuotes(stockFigis)
	})
	wg.Go(func() {
		bondQuotes, err = GetCachedBondQuotes(bondFigis)
	})

	wg.Wait()

	return stockQuotes, bondQuotes, err
}

func GetCachedStockQuotes(figis []string) ([]entity.SimpleQuote, error) {
	stockQuoteMx.Lock()
	defer stockQuoteMx.Unlock()

	targetQuotes := []entity.SimpleQuote{}

	for i := range figis {
		foundQuote := false
		for j := range cachedStockQuotes {
			if figis[i] == cachedStockQuotes[j].Figi() {
				foundQuote = true
				targetQuotes = append(targetQuotes, cachedStockQuotes[j])
				break
			}
		}
		if !foundQuote {
			return targetQuotes, errors.New("Failed to find a cached stock quote for " + figis[i])
		}
	}

	return targetQuotes, nil
}

func GetCachedBondQuotes(tickers []string) ([]entity.BondQuote, error) {
	bondQuoteMx.Lock()
	defer bondQuoteMx.Unlock()

	targetQuotes := []entity.BondQuote{}

	for i := range tickers {
		foundQuote := false
		for j := range cachedBondQuotes {
			if tickers[i] == cachedBondQuotes[j].GetTicker() {
				foundQuote = true
				targetQuotes = append(targetQuotes, cachedBondQuotes[j])
				break
			}
		}
		if !foundQuote {
			return targetQuotes, errors.New("Failed to find a cached bond quote for " + tickers[i])
		}
	}

	return targetQuotes, nil
}
