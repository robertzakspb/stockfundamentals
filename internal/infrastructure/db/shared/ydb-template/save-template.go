package ydbtemplate

import (
	"context"
	"errors"
	"path"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

func SaveEntity(entity types.Value, tablePath string) error {
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		return err
	}
	defer db.ReleaseDriver(dbConnection)

	tablePath = path.Join(dbConnection.Name(), tablePath)

	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tablePath,
		table.BulkUpsertDataRows(entity),
	)

	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save the entity to the database. Reason: " + err.Error())
	}

	return nil
}
