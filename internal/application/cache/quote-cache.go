package cache

import (
	"errors"
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

	quotes, err := quoteservice.FetchBondQuotesFromTapi(bondservice.ExtractBondFigis(&bondList))

	cachedBondQuotes = quotes

	return nil
}
