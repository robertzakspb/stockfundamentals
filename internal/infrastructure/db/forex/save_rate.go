package forexdb

import (
	"context"
	"errors"
	"path"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func SaveForexRates(rates []ForexRateDb) error {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return err
	}
	defer db.Close(context.TODO())

	ydbRates := []types.Value{}
	for _, r := range rates {
		ydbRate := types.StructValue(
			types.StructFieldValue("currency_1", types.TextValue(r.Currency1)),
			types.StructFieldValue("currency_2", types.TextValue(r.Currency2)),
			types.StructFieldValue("date", shared.ConvertToYdbDate(r.Date)),
			types.StructFieldValue("rate", types.DoubleValue(r.Rate)),
		)
		ydbRates = append(ydbRates, ydbRate)
	}

	tableName := path.Join(db.Name(), shared.FOREX_DIRECTORY_PREFIX, shared.FX_RATE_TABLE_NAME)
	err = db.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbRates...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save forex rates to the database")
	}

	return nil
}
