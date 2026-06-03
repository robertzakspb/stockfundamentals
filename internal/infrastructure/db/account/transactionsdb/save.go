package transactionsdb

import (
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
)

func SaveTransactions(dbModels []TransactionDbModel) error {
	mappedYdbList := mapToYdbList(dbModels)
	tablePath := ydbhelper.GenerateTablePath(db.USER_DIRECTORY_PREFIX, db.TRANSACTION_TABLE_NAME)

	err := ydbtemplate.SaveEntity(mappedYdbList, tablePath)

	return err
}
