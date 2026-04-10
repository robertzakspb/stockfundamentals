package accountmvdb

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

func SaveMarketValue(marketValues []AccountMarketValueDB) error {
	dbConnection, err := db.MakeYdbDriver()
	if err != nil {
		return err
	}
	defer dbConnection.Close(context.TODO())

	dbMarketValues := []types.Value{}
	for _, mv := range marketValues {
		dbMarketValue := types.StructValue(
			types.StructFieldValue("account_id", types.UuidValue(mv.AccountId)),
			types.StructFieldValue("date", ydbhelper.ConvertToYdbDate(mv.Date)),
			types.StructFieldValue("currency", types.TextValue(mv.Currency)),
			types.StructFieldValue("eod_value", types.DoubleValue(mv.EodValue)),
		)
		dbMarketValues = append(dbMarketValues, dbMarketValue)
	}

	tableName := path.Join(dbConnection.Name(), db.USER_DIRECTORY_PREFIX, db.ACCOUNT_MARKET_VALUE_HISTORY_TABLE_NAME)
	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(dbMarketValues...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save bond position lots to the database")
	}

	return nil
}
