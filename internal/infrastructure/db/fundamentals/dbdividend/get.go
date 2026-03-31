package dbdividend

import (
	"context"
	"errors"

	"io"

	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

func GetFilteredDividends(filters []ydbfilter.YdbFilter) ([]dividend.Dividend, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []dividend.Dividend{}, err
	}
	defer db.Close(context.TODO())

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
		YqlColumnName:  "record_date", //TODO: refactor to pull the value dynamically using reflect
		Condition:      ydbfilter.GreaterThanOrEqualTo,
		ConditionValue: shared.ConvertToYdbDate(time.Now()),
	}
	allDividends, err := GetFilteredDividends([]ydbfilter.YdbFilter{payoutDateFilter})
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

func GetDividendForecasts() ([]DividendForecastDb, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []DividendForecastDb{}, err
	}
	defer db.Close(context.TODO())

	forecasts := []DividendForecastDb{}

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetDividendForecastQuery(),
				query.WithTxControl(query.TxControl(query.BeginTx(query.WithSnapshotReadOnly()))))

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

				for row, err := range sugar.UnmarshalRows[DividendForecastDb](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					forecasts = append(forecasts, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		return []DividendForecastDb{}, err
	}

	return forecasts, nil
}
