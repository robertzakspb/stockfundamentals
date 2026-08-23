package transactionsdb

import (
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
)

func GetAllTransactions() ([]TransactionDbModel, error) {
	tablePath := "`" + ydbhelper.GenerateTablePath(db.USER_DIRECTORY_PREFIX, db.TRANSACTION_TABLE_NAME) + "`"
	filters := []ydbfilter.YdbFilter{}

	transactions, err := ydbtemplate.GetEntity[TransactionDbModel](filters, tablePath)
	return transactions, err
}

func GetFilteredTransactions(query shared.ParsedApiQuery) ([]TransactionDbModel, error) {
	tablePath := "`" + ydbhelper.GenerateTablePath(db.USER_DIRECTORY_PREFIX, db.TRANSACTION_TABLE_NAME) + "`"

	transactions, err := ydbtemplate.GetFilteredEntity[TransactionDbModel](query.Filters, query, tablePath)
	return transactions, err
}
