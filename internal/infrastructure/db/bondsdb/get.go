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
							coupon_date,
							coupon_number,
							record_date,
							per_bond_amount,
							coupon_type,
							coupon_start_date,
							coupon_end_date
						FROM
							%s
					`,
		"`"+path.Join(shared.BOND_DIRECTORY_PREFIX, shared.BOND_TABLE_NAME)+"`")

	return yql
}
