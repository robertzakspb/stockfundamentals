package transactionsdb

import (
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
	"github.com/google/uuid"
)

func DeleteTbankTransactions(accountIds []uuid.UUID) error {
	filters := []ydbfilter.YdbFilter{
		{
			YqlColumnName:  "account_id",
			Condition:      ydbfilter.Contains,
			ConditionValue: ydbhelper.ConvertUUIDsToYdbList(accountIds),
		},
	}
	tablePath := "`" + ydbhelper.GenerateTablePath(db.USER_DIRECTORY_PREFIX, db.TRANSACTION_TABLE_NAME) + "`"

	err := ydbtemplate.DeleteEntity(filters, tablePath)

	return err
}
