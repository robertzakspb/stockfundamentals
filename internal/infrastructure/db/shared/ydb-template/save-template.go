package ydbtemplate

import (
	"context"
	"errors"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func SaveEntity(entities []types.Value, tablePath string) error {
	if len(entities) == 0 {
		return errors.New("Attempting to save 0 entities to the database")
	}

	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		return err
	}
	defer db.ReleaseDriver(dbConnection)

	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tablePath,
		table.BulkUpsertDataRows(types.ListValue(entities...)),
	)

	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save bond position lots to the database")
	}

	return nil
}