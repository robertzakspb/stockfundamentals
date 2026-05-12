package timeseries

import (
	"errors"

	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/quote"
	timeseriesdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
)

// This function will attempt to fetch the latest quotes from the database and, if any are missing, attempt to fetch them from supported 3rd parties
func GetLatestQuotes(figis []string) ([]quote.Quote, error) {
	quotes := []quote.Quote{}

	dbQuotes, err := timeseriesdb.GetLatestQuotesForAllSecurities()
	if err != nil {
		return quotes, err
	}

	quotes = append(quotes, mapDbQuotesToQuotes(dbQuotes)...)

	tickersWithMissingQuotes := []string{}
	for i := range figis {
		foundQuote := false
		for j := range dbQuotes {
			if figis[i] == dbQuotes[j].Figi {
				foundQuote = true
			}
		}
		if !foundQuote {
			tickersWithMissingQuotes = append(tickersWithMissingQuotes, figis[i])
		}
	}

	missingQuotes, err := FetchTinkoffQuotesFor(tickersWithMissingQuotes)
	if err != nil {
		return quotes, err
	}

	quotes = append(quotes, missingQuotes...)

	return quotes, nil
}

func FetchTinkoffQuotesFor(figis []string) ([]quote.Quote, error) {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		return []quote.Quote{}, errors.New("Unable to fetch quotes due to internal configuration issues")
	}

	quotes, err := tquoteservice.FetchQuotesForFigis(figis, config)
	if err != nil {
		return []quote.Quote{}, errors.New(err.Error())
	}

	mappedQuotes := mapTQuotesToQuotes(quotes)

	return mappedQuotes, nil
}
