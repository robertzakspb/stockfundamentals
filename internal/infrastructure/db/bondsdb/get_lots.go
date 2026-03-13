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
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

func GetAccountBondPortfolio(filters []ydbfilter.YdbFilter) ([]BondPositionLotDb, error) {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return []BondPositionLotDb{}, err
	}

	bonds := []BondPositionLotDb{}

	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx, makeGetBondPositionsQuery(filters),
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

				for row, err := range sugar.UnmarshalRows[BondPositionLotDb](resultSet.Rows(ctx)) {
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
		logger.Log(err.Error(), logger.ERROR)
		return []BondPositionLotDb{}, err
	}

	return bonds, nil
}

func makeGetBondPositionsQuery(filters []ydbfilter.YdbFilter) string {
	yql := fmt.Sprintf(`
						%s
						SELECT
							id,
							figi,
							isin,
							opening_date,
							modification_date,
							account_id,
							quantity,
							price_per_unit
						FROM
							%s
						%s
					`,
		ydbfilter.AddYqlVarDeclarations(filters),
		"`"+path.Join(shared.BOND_DIRECTORY_PREFIX, shared.BOND_POSITION_LOT_TABLE_NAME)+"`",
		ydbfilter.MakeWhereClause(filters))
	return yql
}
