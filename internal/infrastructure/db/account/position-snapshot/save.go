package positionsnapshotdb

import (
	"path"

	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
)

func SaveStockPositionSnapshots(dbModels []StockPositionSnapshotDbModel) error {
	tablePath := path.Join(db.USER_DIRECTORY_PREFIX, db.STOCK_LOT_SNAPSHOT_TABLE_NAME)
	ydbList := mapSnapshotDbModelToYdbList(dbModels)

	err := ydbtemplate.SaveEntity(ydbList, tablePath)

	return err
}
