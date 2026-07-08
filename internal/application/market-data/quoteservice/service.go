package quoteservice

import (
	"context"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

func FetchBondQuotes(tickers []string) ([]entity.BondQuote, error) {
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
	batches := stringhelpers.SplitInBatchesOf(100, tickers)
	bondQuotes := []entity.BondQuote{}

	for i := range batches {
		quotes, errorList := tquoteservice.GetBondPriceAndYield(client, batches[i])
		if len(errorList) > 0 {
			logger.Log("Failed to fetch bond quotes due to: "+errorList[0].Error(), logger.ERROR)
			return []entity.BondQuote{}, err
		}
		bondQuotes = append(bondQuotes, quotes...)
	}

	return bondQuotes, nil
}

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
