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

func GetBondCoupons(filters []ydbfilter.YdbFilter) ([]CouponDbModel, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []CouponDbModel{}, err
	}
	defer db.Close(context.TODO())

	coupons := []CouponDbModel{}

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetCouponsByFigiQuery(filters),
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

				for row, err := range sugar.UnmarshalRows[CouponDbModel](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					coupons = append(coupons, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		return []CouponDbModel{}, err
	}

	return coupons, nil
}

func makeGetCouponsByFigiQuery(filters []ydbfilter.YdbFilter) string {
	yql := fmt.Sprintf(`
						%s
						SELECT
							id,
							figi,
							coupon_date,
							coupon_number,
							record_date,
							per_bond_amount,
							coupon_type,
							coupon_start_date,
							coupon_end_date,
							coupon_period
						FROM
							%s
						%s
					`,
		ydbfilter.AddYqlVarDeclarations(filters),
		"`"+path.Join(shared.BOND_DIRECTORY_PREFIX, shared.COUPON_TABLE_NAME)+"`",
		ydbfilter.MakeWhereClause(filters))

	return yql
}
