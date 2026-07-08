package timeseries

import (
	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/quote"
	timeseriesdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func mapTQuotesToQuotes(tQuotes []tquoteservice.TQuote) []quote.Quote {
	quotes := make([]quote.Quote, len(tQuotes))

	for i := range tQuotes {
		mappedQuote := quote.New(tQuotes[i].Figi(), tQuotes[i].Currency(), tQuotes[i].Timestamp(), tQuotes[i].Quote())
		quotes[i] = mappedQuote
	}

	return quotes
}

func mapDbQuotesToQuotes(dbQuotes []timeseriesdb.QuoteDB) []quote.Quote {
	quotes := make([]quote.Quote, len(dbQuotes))

	for i := range dbQuotes {
		currency := ""
		switch dbQuotes[i].Country {
		case "RU":
			currency = "RUB"
		case "US":
			currency = "USD"
		case "RS":
			currency = "RSD"
		default:
			logger.Log("Failed to map the country "+dbQuotes[i].Country+" to a corresponding currency", logger.ERROR)
			currency = dbQuotes[i].Country
		}

		mappedQuoted := quote.New(dbQuotes[i].Figi, currency, dbQuotes[i].Date, dbQuotes[i].ClosePrice)
		quotes[i] = mappedQuoted
	}

	return quotes
}

func mapDbQuoteToQuote(dbQuote timeseriesdb.QuoteDB) quote.Quote {
	quote := mapDbQuotesToQuotes([]timeseriesdb.QuoteDB{dbQuote})[0]
	return quote
}
