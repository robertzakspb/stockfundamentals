package forexdb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"slices"
	"strings"
	"time"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

// Strict function that returns forex rates with the exact filters applied. If you're satisfied with any rate for a currency in the last 7 days, use the
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

// Returns the latest available forex rates in the DB given the filters. Do not pass any date filters in the arguments
func GetLatestAvailableRates(filters []ydbfilter.YdbFilter) ([]ForexRateDb, error) {
	//To get the latest forex rates, we fetch the rates for the last 7 days and get the latest available in the selection
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	dateRangeFilter := []ydbfilter.YdbFilter{
		{
			YqlColumnName:  "date",
			Condition:      ydbfilter.GreaterThanOrEqualTo,
			ConditionValue: ydbhelper.ConvertToYdbDate(sevenDaysAgo),
		},
		{
			YqlColumnName:  "date",
			Condition:      ydbfilter.LessThanOrEqualTo,
			ConditionValue: ydbhelper.ConvertToYdbDate(time.Now()),
		},
	}
	filters = append(filters, dateRangeFilter...)

	tablePath := "`" + path.Join(db.FOREX_DIRECTORY_PREFIX, db.FX_RATE_TABLE_NAME) + "`"
	rates, err := ydbtemplate.GetFilteredEntity[ForexRateDb](shared.ParsedApiQuery{Filters: filters}, tablePath)

	//The fetched rates contain the rates for the last 7 days, we need to find the latest one for each pair and return it
	ratesGroupedByCurrency := map[string][]ForexRateDb{}
	for i := range rates {
		pair := strings.Join([]string{rates[i].Currency1, rates[i].Currency2}, "/")
		ratesGroupedByCurrency[pair] = append(ratesGroupedByCurrency[pair], rates[i])
	}
	latestRates := []ForexRateDb{}
	for _, rates := range ratesGroupedByCurrency {
		slices.SortFunc(rates, func(r1, r2 ForexRateDb) int {
			if timehelpers.DateIsLater(r1.Date, r2.Date) {
				return 1
			} else if timehelpers.DateIsEarlier(r1.Date, r2.Date) {
				return -1
			} else {
				return 0
			}
		})
		latestRates = append(latestRates, rates[len(rates)-1])
	}

	return latestRates, err
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
		"`"+path.Join(db.FOREX_DIRECTORY_PREFIX, db.FX_RATE_TABLE_NAME)+"`",
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
		"`"+path.Join(db.FOREX_DIRECTORY_PREFIX, db.FX_RATE_TABLE_NAME)+"`",
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

				}
			}

			return nil
		},
	)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return rates.MinDate, rates.MaxDate, nil
}
