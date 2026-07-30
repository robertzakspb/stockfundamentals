package quoteservice

import (
	"context"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/timeseries"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

// Only fetches bond quotes from the T Bank API
func FetchBondQuotes(figis []string) ([]entity.BondQuote, error) {
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
			//The error "The YTM was not found in the response" is actually not that critical and can be just a warning
			for i := range errorList {
				if errorList[i].Error() == "The YTM was not found in the response" {
					logger.Log(errorList[i].Error(), logger.WARNING)
				} else {
					logger.Log(errorList[i].Error(), logger.ERROR)
				}
			}
			logger.Log("Failed to fetch bond quotes due to: "+errorList[0].Error(), logger.ERROR)
		}
		bondQuotes = append(bondQuotes, quotes...)
	}

	return bondQuotes, nil
}

// Only fetches stock quotes from the T Bank API
func FetchStockQuotes(figis []string) ([]entity.SimpleQuote, error) {
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

	missingQuotes, err := FetchStockQuotes(figisWithMissingQuotes)

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

	missingQuotes, err := FetchBondQuotes(figisWithMissingQuotes)

	quotes = append(quotes, missingQuotes...)

	return quotes, nil
}
