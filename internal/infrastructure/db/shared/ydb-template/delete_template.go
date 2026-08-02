package ydbtemplate

import (
	"context"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

// Deletes all records from the given table
func DeleteEntity(filters []ydbfilter.YdbFilter, tablePath string) error {
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		return err
	}
	defer db.ReleaseDriver(dbConnection)

	yql := generateDeleteQuery(filters, tablePath)

	params := ydbfilter.SetQueryParams(filters)

	err = dbConnection.Table().DoTx(context.TODO(),
		func(ctx context.Context, tx table.TransactionActor) (err error) {
			result, err := tx.Execute(ctx,
				yql,
				params,
			)
			if err != nil {
				return err
			}

			defer func() {
				_ = result.Close()
			}()
			return nil
		})
	if err != nil {
		return err
	}

	return nil
}
