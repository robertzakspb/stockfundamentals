package accountservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	accountdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/account"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
)

func GetAllAccounts() ([]account.Account, error) {
	dbAccounts, err := accountdb.GetAccounts([]ydbfilter.YdbFilter{})
	if err != nil {
		return []account.Account{}, err
	}

	accounts := mapDbAccountsToAccounts(dbAccounts)

	return accounts, nil
}
