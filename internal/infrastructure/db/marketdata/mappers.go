package timeseriesdb

import (
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

func mapBondQuoteDbModelToYdbEntity(dbModels []BondQuoteDB) types.Value {
	dbQuotes := make([]types.Value, len(dbModels))
	for i, quote := range dbModels {
		dbQuote := types.StructValue(
			types.StructFieldValue("ticker", types.TextValue(quote.Ticker)),
			types.StructFieldValue("timestamp", ydbhelper.ConvertToYdbDateTime(quote.Timestamp)),
			types.StructFieldValue("price_as_percentage", types.DoubleValue(quote.PriceAsPercentage)),
		)
		dbQuotes[i] = dbQuote
	}

	return types.ListValue(dbQuotes...)

}
