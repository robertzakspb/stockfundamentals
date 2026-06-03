package accountdb

import (
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
)

func SaveAccountsToDb(dbModels []AccountDbModel) error {
	mappedAccounts := mapAccountDbModelToYdbEntity(dbModels)
	tablePath := ydbhelper.GenerateTablePath(db.USER_DIRECTORY_PREFIX, db.ACCOUNT_TABLE_NAME)

	err := ydbtemplate.SaveEntity(mappedAccounts, tablePath)
	return err
}
