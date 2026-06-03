package accountservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	accountdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/account"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
)

func GetAllAccounts() ([]account.Account, error) {
	dbAccounts, err := accountdb.GetAccounts([]ydbfilter.YdbFilter{})
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return []account.Account{}, err
	}

	accounts := mapDbAccountsToAccounts(dbAccounts)

	return accounts, nil
}

func GetAccountsById(ids uuid.UUIDs) ([]account.Account, error) {
	idFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Contains,
		ConditionValue: ydbhelper.ConvertUUIDsToYdbList(ids),
	}

	dbAccounts, err := accountdb.GetAccounts([]ydbfilter.YdbFilter{idFilter})
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return []account.Account{}, err
	}

	accounts := mapDbAccountsToAccounts(dbAccounts)

	return accounts, nil
}

func SaveAccounts(accounts []account.Account) error {
	dbModels := mapAccountsToDbAccounts(accounts)

	err := accountdb.SaveAccountsToDb(dbModels)

	return err
}
