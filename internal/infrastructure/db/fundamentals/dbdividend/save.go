package dbdividend

import (
	"context"
	"errors"
	"path"
	"strconv"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func SaveDividendsToDB(dividends *[]dividend.Dividend) error {
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}
	defer db.ReleaseDriver(dbConnection)

	if len(*dividends) == 0 {
		logger.Log("Attempting to save 0 dividends", logger.WARNING)
	}

	dbModels := mapDividendToDbModel(*dividends)

	ydbDividends := make([]types.Value, len(*dividends))
	for i, dividend := range dbModels {
		ydbDividend := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(dividend.Id)),
			types.StructFieldValue("stock_id", types.TextValue(dividend.Figi)),
			types.StructFieldValue("actual_DPS", types.Int64Value(int64(dividend.ActualDPSTimesMillion))),
			types.StructFieldValue("expected_DPS", types.Int64Value(int64(dividend.ExpectedDpsTimesMillion))),
			types.StructFieldValue("currency", types.TextValue(dividend.Currency)),
			types.StructFieldValue("announcement_date", ydbhelper.ConvertToOptionalYDBdate(dividend.AnnouncementDate)),
			types.StructFieldValue("record_date", ydbhelper.ConvertToOptionalYDBdate(dividend.RecordDate)),
			types.StructFieldValue("payout_date", ydbhelper.ConvertToOptionalYDBdate(dividend.PayoutDate)),
			types.StructFieldValue("payment_period", types.TextValue(dividend.PaymentPeriod)),
			types.StructFieldValue("management_comment", types.TextValue(dividend.ManagementComment)),
		)
		ydbDividends[i] = ydbDividend
	}

	tableName := path.Join(dbConnection.Name(), db.STOCK_DIRECTORY_PREFIX, db.DIVIDEND_PAYMENT_TABLE_NAME)
	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbDividends...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save dividends to the database")
	}

	logger.Log("Saved "+strconv.Itoa(len(ydbDividends))+" dividends to the database", logger.INFORMATION)

	return nil
}

func SaveDividendForecastToDb(forecast DividendForecastDb) error {
	dbConnection, err := db.MakeYdbDriver()
	if err != nil {
		return err
	}
	defer dbConnection.Close(context.TODO())

	ydbForecast := types.StructValue(
		types.StructFieldValue("id", types.UuidValue(uuid.New())),
		types.StructFieldValue("figi", types.TextValue(forecast.Figi)),
		types.StructFieldValue("expected_DPS", types.DoubleValue(forecast.ExpectedDPS)),
		types.StructFieldValue("currency", types.TextValue(forecast.Currency)),
		types.StructFieldValue("payment_period", types.TextValue(forecast.PaymentPeriod)),
		types.StructFieldValue("comment", types.TextValue(forecast.Comment)),
		types.StructFieldValue("forecast_author", types.TextValue(forecast.Author)),
	)

	tableName := path.Join(dbConnection.Name(), db.STOCK_DIRECTORY_PREFIX, db.DIVIDEND_FORECAST_TABLE_NAME)
	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbForecast)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save dividends to the database")
	}
	return nil
}
