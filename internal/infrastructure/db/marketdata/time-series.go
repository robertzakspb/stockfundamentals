package timeseries

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const market_date_directory_prefix = "marketdata/"
const time_series_table_name = "time_series"

func SaveTimeSeriesToDB(quotes []entity.SimpleQuote) error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := ydb.Open(context.TODO(), config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}

	ydbQuotes := []types.Value{}
	for _, quote := range quotes {
		ydbQuote := types.StructValue(
			types.StructFieldValue("figi", types.TextValue(quote.Figi())),
			types.StructFieldValue("close_price", types.DoubleValue(quote.Quote())),
			types.StructFieldValue("date", utilities.ConvertToYdbDate(quote.Timestamp())),
		)

		ydbQuotes = append(ydbQuotes, ydbQuote)
	}

	tableName := makeTimeSeriesTablePath(db.Name())
	err = db.Table().BulkUpsert(
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
	config, err := config.LoadConfig()
	if err != nil {
		return []QuoteDB{}, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*15)
	defer cancel()
	db, err := ydb.Open(ctx, config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}

	dbQuotes := []QuoteDB{}
	yqlQuery := fmt.Sprintf(
		"SELECT "+
			"`marketdata/time_series`.figi AS figi, MAX(date) AS date, `stockfundamentals/stocks/stock`.country_iso2 AS country_iso2"+
			" FROM "+
			"%s"+
			" JOIN `stockfundamentals/stocks/stock` ON `marketdata/time_series`.figi = `stockfundamentals/stocks/stock`.figi"+
			" GROUP BY `marketdata/time_series`.figi, `stockfundamentals/stocks/stock`.country_iso2",
		makeTimeSeriesTablePath(db.Name()))
	logger.Log("Executing query: "+yqlQuery, logger.INFORMATION)

	err = db.Query().Do(context.TODO(),
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

func makeTimeSeriesTablePath(dbName string) string {
	path := "`" + path.Join( market_date_directory_prefix, time_series_table_name) +  "`"
	return path
}
