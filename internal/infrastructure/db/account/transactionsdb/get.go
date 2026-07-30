package transactionsdb

import (
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
)

func GetAllTransactions() ([]TransactionDbModel, error) {
	tablePath := "`" + ydbhelper.GenerateTablePath(db.USER_DIRECTORY_PREFIX, db.TRANSACTION_TABLE_NAME) + "`"
	filters := []ydbfilter.YdbFilter{}

	transactions, err := ydbtemplate.GetEntity[TransactionDbModel](filters, tablePath)
	return transactions, err
}
