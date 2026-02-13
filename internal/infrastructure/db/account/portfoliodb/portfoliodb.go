package portfoliodb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
)

func GetAccountPortfolio(accountIDs uuid.UUIDs) ([]LotDb, error) {
	db, err := shared.MakeYdbDriver()
	if err != nil {
		return []LotDb{}, err
	}

	lots := []LotDb{}

	//SUGGESTION FROM REDDIT:
	/*
			query := `
		SELECT *
		FROM ` + "`Role`"

	*/

	//FIXME: Complete this code
	/*
			Functioning query:
			SELECT
		    `stockfundamentals/stocks/stock`.figi AS figi,
		    `stockfundamentals/stocks/stock`.company_name AS company_name,
		    `stockfundamentals/stocks/stock`.ticker AS ticker,
		    `user/position_lot`.account_id AS account_id,
		    `user/position_lot`.created_at AS created_at,
		    `user/position_lot`.quantity AS quantity,
		    `user/position_lot`.price_per_unit AS price_per_unit,
		    `user/position_lot`.currency AS currency
		FROM `user/position_lot`
		JOIN `stockfundamentals/stocks/stock` ON `user/position_lot`.figi = `stockfundamentals/stocks/stock`.figi
	*/
	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx,
				makeGetAccountPortfolioQuery(),
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

				for row, err := range sugar.UnmarshalRows[LotDb](resultSet.Rows(ctx)) {
					if err != nil {
						return err
					}

					lots = append(lots, row)
				}
			}

			return nil
		},
	)
	if err != nil {
		fmt.Println(err)
		return []LotDb{}, err
	}

	return lots, nil
}

func UpdateLocalPortfolio(lots []LotDb) error {
	return nil
}

func makeGetAccountPortfolioQuery() string {
	yql := "SELECT" +
		"`stockfundamentals/stocks/stock`.figi AS figi," +
		"`stockfundamentals/stocks/stock`.company_name AS company_name," +
		"`stockfundamentals/stocks/stock`.ticker AS ticker," +
		"`user/position_lot`.account_id AS account_id," +
		"`user/position_lot`.created_at AS created_at," +
		"`user/position_lot`.quantity AS quantity," +
		"`user/position_lot`.currency AS currency, " +
		"`user/position_lot`.price_per_unit AS price_per_unit" +
		" FROM" +
		"`" + path.Join(shared.USER_DIRECTORY_PREFIX, shared.POSITION_LOT_TABLE_NAME) + "`" +
		" JOIN " +
		"`" + path.Join(shared.STOCK_DIRECTORY_PREFIX, shared.STOCK_TABLE_NAME) + "`" +
		" ON " +
		"`" + path.Join(shared.STOCK_DIRECTORY_PREFIX, shared.STOCK_TABLE_NAME) + "`.figi" + " = " +
		"`" + path.Join(shared.USER_DIRECTORY_PREFIX, shared.POSITION_LOT_TABLE_NAME) + "`" + ".figi"
	return yql
}
