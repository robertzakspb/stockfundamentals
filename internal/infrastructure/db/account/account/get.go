package accountdb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

func GetAccounts(filters []ydbfilter.YdbFilter) ([]AccountDbModel, error) {
	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		return []AccountDbModel{}, err
	}
	defer db.ReleaseDriver(dbConnection)

	accounts := []AccountDbModel{}

	err = dbConnection.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetAccountsQuery(filters),
				query.WithTxControl(query.TxControl(query.BeginTx(query.WithSnapshotReadOnly()))),
				query.WithParameters(ydbfilter.SetQueryParams(filters)))

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

				for row, err := range sugar.UnmarshalRows[AccountDbModel](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}
					accounts = append(accounts, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		return []AccountDbModel{}, err
	}

	return accounts, nil
}

func makeGetAccountsQuery(filters []ydbfilter.YdbFilter) string {
	yql := fmt.Sprintf(`
						%s
						SELECT
							id,
							opening_date,
							type,
							broker,
							holder,
							primary_currency
						FROM
							%s
						%s
					`,
		ydbfilter.AddYqlVarDeclarations(filters),
		"`"+path.Join(db.USER_DIRECTORY_PREFIX, db.ACCOUNT_TABLE_NAME)+"`",
		ydbfilter.MakeWhereClause(filters))
	return yql
}
