package bondsdb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

func GetAllBonds(filters []ydbfilter.YdbFilter) ([]BondDbModel, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []BondDbModel{}, err
	}
	defer db.Close(context.TODO())

	bonds := []BondDbModel{}

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetAllBondsQuery(filters),
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

				for row, err := range sugar.UnmarshalRows[BondDbModel](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					bonds = append(bonds, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		return []BondDbModel{}, err
	}

	return bonds, nil
}

func makeGetAllBondsQuery(filters []ydbfilter.YdbFilter) string {
	yql := fmt.Sprintf(`
						%s
						SELECT
							id,
							figi,
							isin,
							lot,
							currency,
							name,
							country_of_risk,
							real_exchange,
							coupon_count_per_year,
							maturity_date,
							nominal_value,
							nominal_currency,
							initial_nominal_value,
							initial_nominal_currency,
							registration_date,
							placement_date,
							placement_price,
							placement_currency,
							accumulated_coupon_income,
							issue_size,
							issue_size_plan,
							has_floating_coupon,
							is_perpetual,
							has_amortization,
							is_available_for_iis,
							is_for_qualified_investors,
							is_subordinated,
							risk_level,
							bond_type,
							call_option_exercise_date
						FROM
							%s
						%s
					`,
		ydbfilter.AddYqlVarDeclarations(filters),
		"`"+path.Join(shared.BOND_DIRECTORY_PREFIX, shared.BOND_TABLE_NAME)+"`",
		ydbfilter.MakeWhereClause(filters))
	return yql
}
