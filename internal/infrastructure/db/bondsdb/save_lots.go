package bondsdb

import (
	"context"
	"errors"
	"path"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	utilities "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func SaveBondPositionLots(lots []BondPositionLotDb) error {
	db, err := utilities.MakeYdbDriver()
	if err != nil {
		return err
	}

	ydbBondLots := []types.Value{}
	for _, l := range lots {
		ydbBondLot := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(l.Id)),
			types.StructFieldValue("figi", types.TextValue(l.Figi)),
			types.StructFieldValue("isin", types.TextValue(l.Isin)),
			types.StructFieldValue("opening_date", shared.ConvertToYdbDateTime(l.OpeningDate)),
			types.StructFieldValue("modification_date", shared.ConvertToYdbDateTime(l.ModificationDate)),
			types.StructFieldValue("account_id", types.UuidValue(l.AccountId)),
			types.StructFieldValue("quantity", types.DoubleValue(l.Quantity)),
			types.StructFieldValue("price_per_unit", types.DoubleValue(l.PricePerUnit)),
		)
		ydbBondLots = append(ydbBondLots, ydbBondLot)
	}

	tableName := path.Join(db.Name(), shared.BOND_DIRECTORY_PREFIX, shared.BOND_POSITION_LOT_TABLE_NAME)
	err = db.Table().BulkUpsert(
		context.TODO(),
		tableName,
		table.BulkUpsertDataRows(types.ListValue(ydbBondLots...)))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return errors.New("Failed to save bond position lots to the database")
	}

	return nil
}
