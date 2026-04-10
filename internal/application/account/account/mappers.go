package accountservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	accountdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/account"
)

func mapDbAccountsToAccounts(dbAccounts []accountdb.AccountDbModel) []account.Account {
	accounts := []account.Account{}

	for _, dbAccount := range dbAccounts {
		account := account.Account{
			Id:          dbAccount.Id,
			OpeningDate: dbAccount.OpeningDate,
			Type:        string(account.AccountType_Map[dbAccount.Type]),
			Broker:      dbAccount.Broker,
			Holder:      dbAccount.Holder,
		}
		accounts = append(accounts, account)
	}

	return accounts
}
