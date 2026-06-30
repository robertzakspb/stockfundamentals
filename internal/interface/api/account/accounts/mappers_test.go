package accountsapi

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_mapAccountsToDtos(t *testing.T) {
	id := uuid.New()
	date := time.Now()
	acc := account.Account{
		Id:          id,
		OpeningDate: date,
		Type: account.IIS_3,
		Broker: "IBKR",
		Holder: "John Appleseed",
		PrimaryCurrency: "EUR",
		CashBalance: 100,
	}

	dtos := mapAccountsToDtos([]account.Account{acc})

	test.AssertEqual(t, 1, len(dtos))
	test.AssertEqual(t, id, dtos[0].Id)
	test.AssertEqual(t, date, dtos[0].OpeningDate)
	test.AssertEqual(t, "IIS_3", dtos[0].Type)
	test.AssertEqual(t, "EUR", dtos[0].PrimaryCurrency)
	test.AssertEqual(t, 100, dtos[0].CashBalance)
}
