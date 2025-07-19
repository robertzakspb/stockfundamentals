package timeseries

import (
	"context"
	"path"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const market_date_directory_prefix = "marketdata/timeseries"
const time_series_table_name = "time_series"

func saveTimeSeriesToDB(quotes []entity.SimpleQuote) error {
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
			types.StructFieldValue("date", convertToYdbDate(quote.Timestamp())),
		)

		ydbQuotes = append(ydbQuotes, ydbQuote)
	}

	tableName := path.Join(db.Name(), market_date_directory_prefix, time_series_table_name)
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

func convertToYdbDate(date time.Time) types.Value {
	const secondsInADay = 86400
	return types.DateValue(uint32(date.Unix() / secondsInADay))
}
