package accountsapi

import "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"

func mapAccountsToDtos(accounts []account.Account) []AccountDto {
	dtos := []AccountDto{}

	for i := range accounts {
		dto := AccountDto{
			Id:              accounts[i].Id,
			OpeningDate:     accounts[i].OpeningDate,
			Type:            string(accounts[i].Type),
			Broker:          accounts[i].Broker,
			Holder:          accounts[i].Holder,
			PrimaryCurrency: accounts[i].PrimaryCurrency,
			CashBalance:     accounts[i].CashBalance,
		}
		dtos = append(dtos, dto)
	}

	return dtos
}
