package bondsdb

import (
	"context"
	"errors"
	"path"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func SaveBondPositionLots(lots []BondPositionLotDb) error {
	dbConnection, err := db.MakeYdbDriver()
	if err != nil {
		return err
	}
	defer dbConnection.Close(context.TODO())

	ydbBondLots := []types.Value{}
	for _, l := range lots {
		ydbBondLot := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(l.Id)),
			types.StructFieldValue("figi", types.TextValue(l.Figi)),
			types.StructFieldValue("isin", types.TextValue(l.Isin)),
			types.StructFieldValue("opening_date", db.ConvertToYdbDateTime(l.OpeningDate)),
			types.StructFieldValue("modification_date", db.ConvertToYdbDateTime(l.ModificationDate)),
			types.StructFieldValue("account_id", types.UuidValue(l.AccountId)),
			types.StructFieldValue("quantity", types.DoubleValue(l.Quantity)),
			types.StructFieldValue("price_per_unit", types.DoubleValue(l.PricePerUnit)),
			types.StructFieldValue("price_per_unit_rub", types.DoubleValue(l.PricePerUnitInRUB)),
		)
		ydbBondLots = append(ydbBondLots, ydbBondLot)
	}

	tableName := path.Join(dbConnection.Name(), db.BOND_DIRECTORY_PREFIX, db.BOND_POSITION_LOT_TABLE_NAME)
	err = dbConnection.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbBondLots...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save bond position lots to the database")
	}

	return nil
}
