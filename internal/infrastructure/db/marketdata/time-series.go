package timeseriesdb

import (
	"context"
	"errors"
	"io"
	"path"

	"github.com/compoundinvest/invest-core/quote/entity"
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const market_date_directory_prefix = "marketdata/"
const time_series_table_name = "time_series"

func SaveTimeSeriesToDB(quotes *[]entity.SimpleQuote) error {
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}
	defer db.ReleaseDriver(dbConnection)

	ydbQuotes := make([]types.Value, len(*quotes))
	for i := range *quotes {
		ydbQuote := types.StructValue(
			types.StructFieldValue("figi", types.TextValue((*quotes)[i].Figi())),
			types.StructFieldValue("close_price", types.DoubleValue((*quotes)[i].Quote())),
			types.StructFieldValue("date", ydbhelper.ConvertToYdbDate((*quotes)[i].Timestamp())),
		)

		ydbQuotes[i] = ydbQuote
	}

	tableName := path.Join(dbConnection.Name(), market_date_directory_prefix, time_series_table_name)
	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbQuotes...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}

func GetLatestQuotesForAllSecurities() ([]QuoteDB, error) {
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}
	defer db.ReleaseDriver(dbConnection)

	dbQuotes := []QuoteDB{}
	yqlQuery :=
		"$noPriceSelection = SELECT " +
			"`stockfundamentals/stocks/stock`.figi AS figi, " +
			"MAX(date) AS date, " +
			"`stockfundamentals/stocks/stock`.country_iso2 AS country_iso2 " +
			" FROM " +
			"`stockfundamentals/stocks/stock` LEFT JOIN `marketdata/time_series` " +
			"USING (figi) " +
			"GROUP BY `stockfundamentals/stocks/stock`.figi, `stockfundamentals/stocks/stock`.country_iso2;" +

			" SELECT n.figi AS figi, n.date AS date, n.country_iso2 AS country_iso2, t.close_price AS close_price " +
			"FROM " +
			"$noPriceSelection as n " +
			"LEFT JOIN `marketdata/time_series` AS t " +
			" ON t.date = n.date AND t.figi = n.figi" +
			" ORDER BY date DESC"
	logger.Log("Executing query: "+yqlQuery, logger.INFORMATION)

	err = dbConnection.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, yqlQuery)

			if err != nil {
				return err
			}

			defer func() {
				_ = result.Close(ctx)
			}()

			for {
				resultSet, err := result.NextResultSet(ctx)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}

					return err
				}

				for row, err := range sugar.UnmarshalRows[QuoteDB](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					dbQuotes = append(dbQuotes, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return []QuoteDB{}, err
	}

	return dbQuotes, nil
}
