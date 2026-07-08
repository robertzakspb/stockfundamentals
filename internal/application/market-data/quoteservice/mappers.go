package quoteservice

import (
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	timeseriesdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
)

func mapBondQuotesToDbModels(quotes []entity.BondQuote) []timeseriesdb.BondQuoteDB {
	dbModels := make([]timeseriesdb.BondQuoteDB, len(quotes))

	for i := range quotes {
		dbModel := timeseriesdb.BondQuoteDB{
			Ticker:            quotes[i].Ticker(),
			PriceAsPercentage: quotes[i].QuoteAsPercentage(),
			Timestamp:         time.Now(),
		}
		dbModels[i] = dbModel
	}

	return dbModels
}
