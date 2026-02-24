package dbdividend

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

func GetAllDividends(filters []ydbfilter.YdbFilter) ([]dividend.Dividend, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []dividend.Dividend{}, err
	}

	userDividendsDbModels := []dividendDbModel{}

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetDividendQuery(filters),
				query.WithTxControl(query.TxControl(query.BeginTx(query.WithSnapshotReadOnly()))),
				query.WithParameters(ydbfilter.SetQueryParams(filters)),
			)

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

				for row, err := range sugar.UnmarshalRows[dividendDbModel](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					userDividendsDbModels = append(userDividendsDbModels, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		return []dividend.Dividend{}, err
	}

	return mapDbModelToDividend(userDividendsDbModels), nil
}

func GetUpcomingDividends() ([]dividend.Dividend, error) {
	payoutDateFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "payout_date", //TODO: refactor to pull the value dynamically using reflect
		Condition:      ydbfilter.GreaterThanOrEqualTo,
		ConditionValue: shared.ConvertToYdbDate(time.Now()),
	}
	allDividends, err := GetAllDividends([]ydbfilter.YdbFilter{payoutDateFilter})
	if err != nil {
		return []dividend.Dividend{}, err
	}

	upcomingDivs := []dividend.Dividend{}
	for _, div := range allDividends {
		if div.PayoutDate.After(time.Now()) {
			upcomingDivs = append(upcomingDivs, div)
		}
	}

	return upcomingDivs, nil
}

func makeGetDividendQuery(filters []ydbfilter.YdbFilter) string {
	yql := fmt.Sprintf(`
						%s

						SELECT
							id,
							stock_id,
							actual_DPS,
							expected_DPS,
							currency,
							announcement_date,
							record_date,
							payout_date,
							payment_period,
							management_comment
						FROM
							%s
						%s
					`,
		ydbfilter.AddYqlVarDeclarations(filters),
		"`"+path.Join(shared.STOCK_DIRECTORY_PREFIX, shared.DIVIDEND_PAYMENT_TABLE_NAME)+"`",
		ydbfilter.MakeWhereClause(filters))

	fmt.Println(yql)
	return yql
}
