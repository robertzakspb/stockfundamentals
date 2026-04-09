package accountmvdb

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

func GetAccountMarketValues(filters []ydbfilter.YdbFilter) ([]AccountMarketValueDB, error) {
	// filters := []ydbfilter.YdbFilter{}
	// filters = append(filters, ydbfilter.YdbFilter{
	// 	YqlColumnName:  "account_id",
	// 	Condition:      ydbfilter.Equal,
	// 	ConditionValue: types.UuidValue(accountId),
	// })
	// filters = append(filters, ydbfilter.YdbFilter{
	// 	YqlColumnName:  "date",
	// 	Condition:      ydbfilter.Contains,
	// 	ConditionValue: ydbhelper.ConverTimestampsToYdbDates(dates),
	// })

	dbConnection, err := db.GetReusableYdbDriver()
	if err != nil {
		return []AccountMarketValueDB{}, err
	}
	defer db.ReleaseDriver(dbConnection)

	marketValues := []AccountMarketValueDB{}

	err = dbConnection.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetAccountMarketValuesQuery(filters),
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

				for row, err := range sugar.UnmarshalRows[AccountMarketValueDB](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}
					marketValues = append(marketValues, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		return []AccountMarketValueDB{}, err
	}

	return marketValues, nil
}

func makeGetAccountMarketValuesQuery(filters []ydbfilter.YdbFilter) string {
	yql := fmt.Sprintf(`
						%s
						SELECT
							account_id,
							date,
							currency,
							eod_value
						FROM
							%s
						%s
						ORDER BY date ASC
					`,
		ydbfilter.AddYqlVarDeclarations(filters),
		"`"+path.Join(db.USER_DIRECTORY_PREFIX, db.ACCOUNT_MARKET_VALUE_HISTORY_TABLE_NAME)+"`",
		ydbfilter.MakeWhereClause(filters))
	return yql
}
