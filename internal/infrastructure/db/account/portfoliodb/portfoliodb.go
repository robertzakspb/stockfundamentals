package portfoliodb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func GetAccountPortfolio(accountIDs uuid.UUIDs) ([]LotDb, error) {
	db, err := shared.MakeYdbDriver()
	if err != nil {
		return []LotDb{}, err
	}

	lots := []LotDb{}
	yql := makeGetAccountPortfolioQuery()
	err = db.Query().Do(context.TODO(),
		func(ctx context.Context, s query.Session) (err error) {
			result, err := s.Query(ctx,
				yql,
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
		return []LotDb{}, err
	}

	return lots, nil
}

func UpdateLocalPortfolio(lots []LotDb) error {
	err := deleteAllLots()
	if err != nil {
		fmt.Println(err.Error())
	}
	if err != nil {
		return err
	}

	db, err := shared.MakeYdbDriver()
	if err != nil {
		return err
	}

	ydbLots := []types.Value{}
	for _, lot := range lots {
		ydbLot := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(lot.Id)),
			types.StructFieldValue("account_id", types.UuidValue(lot.AccountId)),
			types.StructFieldValue("figi", types.UTF8Value(lot.Figi)),
			types.StructFieldValue("created_at", shared.ConvertToYdbDateTime(lot.CreatedAt)),
			types.StructFieldValue("updated_at", shared.ConvertToYdbDateTime(lot.UpdatedAt)),
			types.StructFieldValue("quantity", types.DoubleValue(lot.Quantity)),
			types.StructFieldValue("price_per_unit", types.DoubleValue(lot.PricePerUnit)),
			types.StructFieldValue("currency", types.UTF8Value(lot.Currency)),
		)
		ydbLots = append(ydbLots, ydbLot)
	}

	lotTableName := path.Join(db.Name(), shared.USER_DIRECTORY_PREFIX, shared.POSITION_LOT_TABLE_NAME)
	err = db.Table().BulkUpsert(
		context.TODO(),
		lotTableName,
		table.BulkUpsertDataRows(types.ListValue(ydbLots...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to update position lots in the database")
	}

	return nil
}

func deleteAllLots() error {
	db, err := shared.MakeYdbDriver()
	if err != nil {
		return err
	}

	yql := "DELETE FROM " + "`" + path.Join(shared.USER_DIRECTORY_PREFIX, shared.POSITION_LOT_TABLE_NAME) + "`"

	//func(ctx context.Context, s Session) error
	err = db.Table().DoTx(context.TODO(),
		func(ctx context.Context, tx table.TransactionActor) (err error) {
			result, err := tx.Execute(ctx,
				yql,
				nil,
			)
			if err != nil {
				return err
			}

			defer func() {
				_ = result.Close()
			}()
			return nil
		})
	if err != nil {
		return err
	}

	return nil
}

func makeGetAccountPortfolioQuery() string {
	yql := "SELECT" +
		"`stockfundamentals/stocks/stock`.figi AS figi," +
		"`stockfundamentals/stocks/stock`.company_name AS company_name," +
		"`stockfundamentals/stocks/stock`.ticker AS ticker," +
		"`user/position_lot`.id AS id," +
		"`user/position_lot`.account_id AS account_id," +
		"`user/position_lot`.created_at AS created_at," +
		"`user/position_lot`.updated_at AS updated_at," +
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
