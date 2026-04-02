package dbsecurity

import (
	"context"
	"errors"
	"io"

	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

func GetFilteredSecurities(filters []ydbfilter.YdbFilter) ([]security.Stock, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []security.Stock{}, err
	}
	defer db.Close(context.TODO())

	stocks := []security.Stock{}

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetFilteredSecuritiesQuery(filters),
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

				for row, err := range sugar.UnmarshalRows[security.Stock](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}
					stocks = append(stocks, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		return []security.Stock{}, err
	}

	return stocks, nil
}

func makeGetFilteredSecuritiesQuery(filters []ydbfilter.YdbFilter) string {
	yqlQuery := getSecuritiesBaseQuery()

	yqlQuery = ydbfilter.AddYqlVarDeclarations(filters) + " " + yqlQuery

	yqlQuery += ydbfilter.MakeWhereClause(filters)

	return yqlQuery
}
