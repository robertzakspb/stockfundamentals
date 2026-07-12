package quoteservice

import (
	"github.com/compoundinvest/invest-core/quote/entity"
	timeseriesdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
)

func MapBondQuotesToDbModels(quotes []entity.BondQuote) []timeseriesdb.BondQuoteDB {
	dbModels := make([]timeseriesdb.BondQuoteDB, len(quotes))

	for i := range quotes {
		dbModel := timeseriesdb.NewBondQuoteDb("", quotes[i].GetTicker(), quotes[i].GetTimestamp(), quotes[i].GetQuoteAsPercentage())
		dbModels[i] = dbModel
	}

	return dbModels
}
