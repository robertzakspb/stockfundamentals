package forexdb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"time"

	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

func GetAllFxRates(filters []ydbfilter.YdbFilter) ([]ForexRateDb, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []ForexRateDb{}, err
	}
	defer db.Close(context.TODO())

	rates := []ForexRateDb{}

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetAllForexRatesQuery(filters),
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

				for row, err := range sugar.UnmarshalRows[ForexRateDb](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					rates = append(rates, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		return []ForexRateDb{}, err
	}

	return rates, nil
}

func makeGetAllForexRatesQuery(filters []ydbfilter.YdbFilter) string {
	yql := fmt.Sprintf(`
						%s
						SELECT
							currency_1,
							currency_2,
							date,
							rate
						FROM
							%s
						%s
					`,
		ydbfilter.AddYqlVarDeclarations(filters),
		"`"+path.Join(shared.FOREX_DIRECTORY_PREFIX, shared.FX_RATE_TABLE_NAME)+"`",
		ydbfilter.MakeWhereClause(filters))
	return yql
}

// Only used in the function below
type minMaxFxRate struct {
	MinDate time.Time `sql:"min_date"`
	MaxDate time.Time `sql:"max_date"`
}

// Returns an array where the first element is the earliest and the second element -- the latest forex rate for the provided currencies in the database
func GetEarliestAndLatestDbRateFor(cur1, cur2 string) (time.Time, time.Time, error) {

	yql := fmt.Sprintf(`
						SELECT
							MIN(date) AS min_date,
							MAX(date) AS max_date,
						FROM
							%s
						WHERE currency_1 = '%s' AND currency_2 = '%s'
					`,
		"`"+path.Join(shared.FOREX_DIRECTORY_PREFIX, shared.FX_RATE_TABLE_NAME)+"`",
		cur1, cur2)

	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	defer db.Close(context.TODO())

	var rates minMaxFxRate

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, yql,
				query.WithTxControl(query.TxControl(query.BeginTx(query.WithSnapshotReadOnly()))),
			)

			if err != nil {
				return err
			}

			defer func() {
				_ = result.Close(ctx)
			}()
			counter := 0
			for {
				resultSet, err := result.NextResultSet(ctx)
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}

					return err
				}

				for row, err := range sugar.UnmarshalRows[minMaxFxRate](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}
					rates = row
					counter += 1
				}
			}

			if counter > 1 {
				logger.Log("Retrieved more than 1 min/max fx rate", logger.ERROR)
			}

			return nil
		},
	)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return rates.MinDate, rates.MaxDate, nil
}
