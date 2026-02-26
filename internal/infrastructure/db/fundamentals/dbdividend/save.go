package dbdividend

import (
	"context"
	"errors"
	"path"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func SaveDividendsToDB(dividends []dividend.Dividend, db *ydb.Driver) error {
	if len(dividends) == 0 {
		logger.Log("Attempting to save 0 dividends", logger.WARNING)
	}
	if db == nil {
		logger.Log("Database driver is nil while attempting to save dividends to the DB", logger.ALERT)
		return errors.New("Database issues")
	}

	dbModels := mapDividendToDbModel(dividends)

	ydbDividends := []types.Value{}
	for _, dividend := range dbModels {
		ydbDividend := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(dividend.Id)),
			types.StructFieldValue("stock_id", types.TextValue(dividend.Figi)),
			types.StructFieldValue("actual_DPS", types.Int64Value(int64(dividend.ActualDPSTimesMillion))),
			types.StructFieldValue("expected_DPS", types.Int64Value(int64(dividend.ExpectedDpsTimesMillion))),
			types.StructFieldValue("currency", types.TextValue(dividend.Currency)),
			types.StructFieldValue("announcement_date", shared.ConvertToOptionalYDBdate(dividend.AnnouncementDate)),
			types.StructFieldValue("record_date", shared.ConvertToOptionalYDBdate(dividend.RecordDate)),
			types.StructFieldValue("payout_date", shared.ConvertToOptionalYDBdate(dividend.PayoutDate)),
			types.StructFieldValue("payment_period", types.TextValue(dividend.PaymentPeriod)),
			types.StructFieldValue("management_comment", types.TextValue(dividend.ManagementComment)),
		)
		ydbDividends = append(ydbDividends, ydbDividend)
	}

	tableName := path.Join(db.Name(), shared.STOCK_DIRECTORY_PREFIX, shared.DIVIDEND_PAYMENT_TABLE_NAME)
	err := db.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbDividends...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save dividends to the database")
	}

	return nil
}
