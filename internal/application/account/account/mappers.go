package accountservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	accountdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/account"
)

func mapDbAccountsToAccounts(dbAccounts []accountdb.AccountDbModel) []account.Account {
	accounts := make([]account.Account, len(dbAccounts))

	for i, dbAccount := range dbAccounts {
		account := account.Account{
			Id:              dbAccount.Id,
			OpeningDate:     dbAccount.OpeningDate,
			Type:            string(account.AccountType_Map[dbAccount.Type]),
			Broker:          dbAccount.Broker,
			Holder:          dbAccount.Holder,
			PrimaryCurrency: dbAccount.PrimaryCurrency,
			CashBalance:     dbAccount.CashBalance,
		}
		accounts[i] = account
	}

	return accounts
}

func mapAccountsToDbAccounts(accounts []account.Account) []accountdb.AccountDbModel {
	dbAccounts := make([]accountdb.AccountDbModel, len(accounts))

	for i, account := range accounts {
		account := accountdb.AccountDbModel{
			Id:              account.Id,
			OpeningDate:     account.OpeningDate,
			Type:            string(account.Type),
			Broker:          account.Broker,
			Holder:          account.Holder,
			PrimaryCurrency: account.PrimaryCurrency,
			CashBalance:     account.CashBalance,
		}
		dbAccounts[i] = account
	}

	return dbAccounts
}
