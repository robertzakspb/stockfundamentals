package accountsapi

import (
	"errors"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/google/uuid"
)

func MapAccountsToDtos(accounts []account.Account) []AccountDto {
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

func mapDtoToAccount(dto NewAccountDto) (account.Account, error) {
	accType, found := account.AccountType_Map[dto.Type]
	if !found {
		return account.Account{}, errors.New("Unknown account type: " + dto.Type)
	}
	account := account.Account{
		Id:              uuid.New(),
		OpeningDate:     dto.OpeningDate,
		Type:            accType,
		Broker:          dto.Broker,
		Holder:          dto.Holder,
		PrimaryCurrency: dto.PrimaryCurrency,
		CashBalance:     dto.CashBalance,
	}
	return account, nil

}
