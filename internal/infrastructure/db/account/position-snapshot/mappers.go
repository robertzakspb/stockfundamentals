package positionsnapshotdb

import (
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

func mapSnapshotDbModelToYdbList(dbModels []StockPositionSnapshotDbModel) types.Value {
	dbSnapshots := make([]types.Value, len(dbModels))

	for i := range dbModels {
		dbSnapshot := types.StructValue(
			types.StructFieldValue("account_id", types.UuidValue(dbModels[i].AccountId)),
			types.StructFieldValue("figi", types.TextValue(dbModels[i].Figi)),
			types.StructFieldValue("isin", types.TextValue(dbModels[i].Isin)),
			types.StructFieldValue("date", ydbhelper.ConvertToYdbDate(dbModels[i].Date)),
			types.StructFieldValue("quantity", types.DoubleValue(dbModels[i].Quantity)),
		)
		dbSnapshots[i] = dbSnapshot
	}

	return types.ListValue(dbSnapshots...)
}
