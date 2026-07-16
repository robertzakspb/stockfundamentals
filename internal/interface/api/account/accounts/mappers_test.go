package accountsapi

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_MapAccountsToDtos_Positive(t *testing.T) {
	id := uuid.New()
	date := time.Now()
	acc := account.Account{
		Id:              id,
		OpeningDate:     date,
		Type:            account.IIS_3,
		Broker:          "IBKR",
		Holder:          "John Appleseed",
		PrimaryCurrency: "EUR",
		CashBalance:     100,
		IsOpen:          true,
	}

	dtos := MapAccountsToDtos([]account.Account{acc})

	test.AssertEqual(t, 1, len(dtos))
	test.AssertEqual(t, id, dtos[0].Id)
	test.AssertEqual(t, date, dtos[0].OpeningDate)
	test.AssertEqual(t, "IIS_3", dtos[0].Type)
	test.AssertEqual(t, "EUR", dtos[0].PrimaryCurrency)
	test.AssertEqual(t, 100, dtos[0].CashBalance)
	test.AssertEqual(t, true, dtos[0].IsOpen)
}

func Test_mapDtoToAccount_Negative_UnknownAccountType(t *testing.T) {
	date := time.Now()
	dto := NewAccountDto{
		OpeningDate:     date,
		Type:            "fake",
		Broker:          "IBKR",
		Holder:          "John Appleseed",
		PrimaryCurrency: "EUR",
		CashBalance:     100,
		IsOpen:          false,
	}

	_, err := mapDtoToAccount(dto)

	test.AssertError(t, err)
}

func Test_mapDtoToAccount_Positive(t *testing.T) {
	date := time.Now()
	broker := "IBKR"
	holder := "John Appleseed"
	currency := "EUR"
	dto := NewAccountDto{
		OpeningDate:     date,
		Type:            "IIS_3",
		Broker:          broker,
		Holder:          holder,
		PrimaryCurrency: currency,
		CashBalance:     100,
		IsOpen:          false,
	}

	mappedAccount, err := mapDtoToAccount(dto)

	test.AssertNoError(t, err)
	test.AssertEqual(t, date, mappedAccount.OpeningDate)
	test.AssertEqual(t, account.IIS_3, mappedAccount.Type)
	test.AssertEqual(t, broker, mappedAccount.Broker)
	test.AssertEqual(t, holder, mappedAccount.Holder)
	test.AssertEqual(t, currency, mappedAccount.PrimaryCurrency)
	test.AssertEqual(t, 100, mappedAccount.CashBalance)
	test.AssertEqual(t, false, mappedAccount.IsOpen)
}
