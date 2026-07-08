package timeseries

import (
	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/quote"
	timeseriesdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
)

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


