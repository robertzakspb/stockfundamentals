package tranlotrelationdb

import (
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type DbModel TransactionLotRelationDb

func mapDbModelsToYdbList(dbModels []TransactionLotRelationDb) types.Value {
	ydbModels := make([]types.Value, len(dbModels))

	for i := range dbModels {
		ydbModel := types.StructValue(
			types.StructFieldValue("id", types.UuidValue(dbModels[i].Id)),
			types.StructFieldValue("transaction_id", types.UuidValue(dbModels[i].TransactionId)),
			types.StructFieldValue("stock_lot_id", types.UuidValue(dbModels[i].StockLotId)),
			types.StructFieldValue("bond_lot_id", types.UuidValue(dbModels[i].BondLotId)),
			types.StructFieldValue("date", ydbhelper.ConvertToYdbDate(dbModels[i].Date)),
			types.StructFieldValue("quantity", types.DoubleValue(dbModels[i].Quantity)),
		)
		ydbModels[i] = ydbModel
	}

	return types.ListValue(ydbModels...)
}
