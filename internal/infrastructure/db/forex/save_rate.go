package forexdb

import (
	"context"
	"errors"
	"path"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func SaveForexRates(rates []ForexRateDb) error {
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		return err
	}
	defer db.ReleaseDriver(dbConnection)

	ydbRates := []types.Value{}
	for _, r := range rates {
		ydbRate := types.StructValue(
			types.StructFieldValue("currency_1", types.TextValue(r.Currency1)),
			types.StructFieldValue("currency_2", types.TextValue(r.Currency2)),
			types.StructFieldValue("date", ydbhelper.ConvertToYdbDate(r.Date)),
			types.StructFieldValue("rate", types.DoubleValue(r.Rate)),
		)
		ydbRates = append(ydbRates, ydbRate)
	}

	tableName := path.Join(dbConnection.Name(), db.FOREX_DIRECTORY_PREFIX, db.FX_RATE_TABLE_NAME)
	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbRates...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save forex rates to the database")
	}

	return nil
}
