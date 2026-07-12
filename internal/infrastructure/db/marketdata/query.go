package timeseriesdb

import (
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
)

func makeLatestBondQuoteQuery() string {
	yql := "DECLARE $ticker_filter1 AS List<Utf8>; "

	yql += `$latestQuote = SELECT ticker, MAX(timestamp) AS timestamp FROM ` +

		db.BackTickPath(db.MARKET_DATA_DIRECTORY_PREFIX, db.BOND_TIME_SERIES_TABLE_NAME) +

		`GROUP BY ticker;

			SELECT t.ticker AS ticker, t.timestamp AS timestamp , q.price_as_percentage AS price_as_percentage, b.figi AS figi 
    		FROM $latestQuote AS t JOIN ` +
		db.BackTickPath(db.MARKET_DATA_DIRECTORY_PREFIX, db.BOND_TIME_SERIES_TABLE_NAME) +
		` AS q ON t.ticker = q.ticker AND t.timestamp = q.timestamp	` +
		` JOIN ` + db.BackTickPath(db.BOND_DIRECTORY_PREFIX, db.BOND_TABLE_NAME) +
		` AS b ON b.ticker = t.ticker`
	yql += " WHERE t.ticker IN $ticker_filter1"
	return yql
}
