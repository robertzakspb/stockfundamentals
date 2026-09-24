package quoteservice

import (
	"context"
	"errors"
	"sync"

	"github.com/compoundinvest/invest-core/quote/belex"
	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/timeseries"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

func FetchInternationalStockAndBondQuotes(securities []entity.Security, bondFigis []string) ([]entity.SimpleQuote, []entity.BondQuote, error) {
	var stockQuotes []entity.SimpleQuote
	var bondQuotes []entity.BondQuote
	var wg sync.WaitGroup
	var err error

	wg.Go(func() {
		stockQuotes, err = FetchInternationalStockQuotes(securities)
	})
	wg.Go(func() {
		bondQuotes, err = FetchBondQuotesFromTapi(bondFigis)
	})
	wg.Wait()

	return stockQuotes, bondQuotes, err
}

func FetchInternationalStockQuotes(securities []entity.Security) ([]entity.SimpleQuote, error) {
	var quotes, tApiQuotes []entity.SimpleQuote
	tApiFigis := []string{}
	belexTickers := []string{}
	belexFigis := []string{}
	var wg sync.WaitGroup
	var err error

	for _, security := range securities {
		switch security.MIC {
		case "MISX":
			tApiFigis = append(tApiFigis, security.Figi)
		case "XBEL":
			belexTickers = append(belexTickers, security.Ticker)
			belexFigis = append(belexFigis, security.Figi)
		default:
			return quotes, errors.New("Unsupported MIC: " + security.MIC)
		}
	}

	wg.Go(func() {
		tApiQuotes, err = FetchStockQuotesFromTapi(tApiFigis)
		quotes = append(quotes, tApiQuotes...)
	})
	wg.Go(func() {
		for i := range belexFigis {
			belexQuote, err := belex.FetchQuoteFor(belexTickers[i], belexFigis[i])
			if err != nil {
				continue
			}
			quotes = append(quotes, belexQuote)
		}
	})
	wg.Wait()

	if err != nil {
		return quotes, err
	}

	return quotes, nil
}

// Fetches the stock and bond quotes in parallel for better performance
func GetTinvestmentStockAndBondQuotes(stockFigis, bondFigis []string) ([]entity.SimpleQuote, []entity.BondQuote, error) {
	var stockQuotes []entity.SimpleQuote
	var bondQuotes []entity.BondQuote
	var err error
	var wg sync.WaitGroup

	wg.Go(func() {
		stockQuotes, err = FetchStockQuotesFromTapi(stockFigis)
	})
	wg.Go(func() {
		bondQuotes, err = FetchBondQuotesFromTapi(stockFigis)
	})
	wg.Wait()

	return stockQuotes, bondQuotes, err
}

// Only fetches bond quotes from the T Bank API
func FetchBondQuotesFromTapi(figis []string) ([]entity.BondQuote, error) {
	config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []entity.BondQuote{}, err
	}
	client, err := investgo.NewClient(context.TODO(), config, nil)
	if err != nil {
		return []entity.BondQuote{}, err
	}

	//T Bank API does not return more than 100 quotes in a given call; hence the batching
	batches := stringhelpers.SplitInBatchesOf(100, figis)
	bondQuotes := []entity.BondQuote{}

	for i := range batches {
		quotes, errorList := tquoteservice.GetBondPriceAndYield(client, batches[i])
		if len(errorList) > 0 {
			for i := range errorList {
				if errorList[i].Error() == "The YTM was not found in the response" || errorList[i].Error() == "The quote was not found in the response" {
					logger.Log(errorList[i].Error(), logger.WARNING)
				} else {
					logger.Log(errorList[i].Error(), logger.ERROR)
				}
			}
		}
		bondQuotes = append(bondQuotes, quotes...)
	}

	return bondQuotes, nil
}

// Only fetches stock quotes from the T Bank API
func FetchStockQuotesFromTapi(figis []string) ([]entity.SimpleQuote, error) {
	config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []entity.SimpleQuote{}, err
	}
	client, err := investgo.NewClient(context.TODO(), config, nil)
	if err != nil {
		return []entity.SimpleQuote{}, err
	}

	//T Bank API does not return more than 100 quotes in a given call; hence the batching
	batches := stringhelpers.SplitInBatchesOf(100, figis)
	stockQuotes := []entity.SimpleQuote{}

	for i := range batches {
		quotes, err := tquoteservice.FetchQuotesForFigis(client, batches[i])
		if err != nil {
			logger.Log("Failed to fetch stock quotes due to: "+err.Error(), logger.ERROR)
			return []entity.SimpleQuote{}, err
		}
		stockQuotes = append(stockQuotes, quotes...)
	}
	return stockQuotes, nil
}

// This function will attempt to fetch the latest quotes from the database and, if any are missing, attempt to fetch them from supported 3rd parties
func GetCachedAndExternalStockQuotes(figis []string) ([]entity.SimpleQuote, error) {
	quotes, err := timeseries.GetLatestLocalQuotesForFigis(figis)
	if err != nil {
		return []entity.SimpleQuote{}, err
	}

	if len(quotes) == len(figis) {
		return quotes, nil
	}

	figisWithMissingQuotes := []string{}
	for i := range figis {
		foundQuote := false
		for j := range quotes {
			if figis[i] == quotes[j].Figi() {
				foundQuote = true
			}
		}
		if !foundQuote {
			figisWithMissingQuotes = append(figisWithMissingQuotes, figis[i])
		}
	}

	missingQuotes, err := FetchStockQuotesFromTapi(figisWithMissingQuotes)

	quotes = append(quotes, missingQuotes...)

	return quotes, err
}

func GetCachedAndExternalBondQuotes(bondList []bonds.Bond) ([]entity.BondQuote, error) {
	tickers := make([]string, len(bondList))
	for i := range bondList {
		tickers[i] = bondList[i].Ticker
	}
	quotes, err := timeseries.GetLatestLocalBondQuotes(tickers)
	if err != nil {
		return quotes, err
	}

	if len(quotes) == len(tickers) {
		return quotes, nil
	}

	figisWithMissingQuotes := []string{}
	for i := range bondList {
		foundQuote := false
		for j := range quotes {
			if bondList[i].Ticker == quotes[j].GetTicker() {
				foundQuote = true
			}
		}
		if !foundQuote {
			figisWithMissingQuotes = append(figisWithMissingQuotes, tickers[i])
		}
	}

	missingQuotes, err := FetchBondQuotesFromTapi(figisWithMissingQuotes)

	quotes = append(quotes, missingQuotes...)	

	return quotes, nil
}
