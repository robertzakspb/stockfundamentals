package timeseries

import (
	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/quoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/quote"
	timeseriesdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
)

// This function will attempt to fetch the latest quotes from the database and, if any are missing, attempt to fetch them from supported 3rd parties
func GetLatestQuotes(figis []string) ([]entity.SimpleQuote, error) {
	quotes := []entity.SimpleQuote{}

	dbQuotes, err := GetLatestLocalQuotesForFigis(figis)
	if err != nil {
		return []entity.SimpleQuote{}, err
	}

	quotes = append(quotes, dbQuotes...)

	tickersWithMissingQuotes := []string{}
	for i := range figis {
		foundQuote := false
		for j := range dbQuotes {
			if figis[i] == dbQuotes[j].Figi() {
				foundQuote = true
			}
		}
		if !foundQuote {
			tickersWithMissingQuotes = append(tickersWithMissingQuotes, figis[i])
		}
	}

	missingQuotes, err := quoteservice.FetchStockQuotes(tickersWithMissingQuotes)
	if err != nil {
		return quotes, err
	}

	quotes = append(quotes, missingQuotes...)

	return quotes, nil
}

// Unlike GetLatestQuotes, this function only retrieves local quotes without calls to 3rd parties
func GetLatestLocalQuotesForAllSecurities() ([]quote.Quote, error) {
	quotes := []quote.Quote{}

	dbQuotes, err := timeseriesdb.GetLatestQuotesForAllSecurities()
	if err != nil {
		return quotes, err
	}

	mappedQuotes := mapDbQuotesToQuotes(dbQuotes)
	quotes = append(quotes, mappedQuotes...)

	return quotes, nil
}

func GetLatestLocalQuotesForFigis(figis []string) ([]entity.SimpleQuote, error) {
	quotes := []entity.SimpleQuote{}

	dbQuotes, err := timeseriesdb.GetLatestQuotesForAllSecurities()
	if err != nil {
		return quotes, err
	}

	for i := range dbQuotes {
		for j := range figis {
			if dbQuotes[i].Figi == figis[j] {
				mappedQuote := mapDbQuoteToQuote(dbQuotes[i])
				quotes = append(quotes, mappedQuote)
			}
		}
	}
	return quotes, nil
}
