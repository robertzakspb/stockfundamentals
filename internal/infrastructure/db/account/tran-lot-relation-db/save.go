package tranlotrelationdb

import (
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
)

func SaveTranLotRelations(dbModels []TransactionLotRelationDb) error {
	ydbList := mapDbModelsToYdbList(dbModels)
	tablePath := ydbhelper.GenerateTablePath(db.USER_DIRECTORY_PREFIX, db.TRANSACTION_LOT_RELATIONSHIP_TABLE_NAME)

	err := ydbtemplate.SaveEntity(ydbList, tablePath)

	return err
}
