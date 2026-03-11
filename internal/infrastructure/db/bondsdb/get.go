package bondsdb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

func GetAllBonds() ([]BondDbModel, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []BondDbModel{}, err
	}

	bonds := []BondDbModel{}

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetAllBondsQuery(),
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

func makeGetAllBondsQuery() string {
	yql := fmt.Sprintf(`
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
					`,
		"`"+path.Join(shared.BOND_DIRECTORY_PREFIX, shared.BOND_TABLE_NAME)+"`")
	return yql
}

func makeGetCouponsByFigiQuery(figi string) string {
	yql := fmt.Sprintf(`
						SELECT
							id,
							figi,
							coupon_date,
							coupon_number,
							record_date,
							per_bond_amount,
							coupon_type,
							coupon_start_date,
							coupon_end_date
						FROM
							%s
						WHERE figi = %s
					`,
		"`"+path.Join(shared.BOND_DIRECTORY_PREFIX, shared.COUPON_TABLE_NAME)+"`", figi)

	return yql
}
