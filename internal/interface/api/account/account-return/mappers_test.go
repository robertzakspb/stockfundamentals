package accountreturnapi

import (
	"testing"
	"time"

	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func test_mapAccountMarketValueToDto_Positive(t *testing.T) {
	accountId := uuid.UUID{}
	currency := "EUR"
	date := time.Now()
	eodValue := 56.5
	mv := accountmvdomain.AccountMarketValue{
		AccountId: accountId,
		Currency:  currency,
		Date:      date,
		EodValue:  eodValue,
	}

	dto := mapAccountMarketValueToDto(mv)

	test.AssertEqual(t, accountId, dto.AccountId)
	test.AssertEqual(t, currency, dto.Currency)
	test.AssertEqual(t, date, dto.Date)
	test.AssertEqual(t, eodValue, dto.EodValue)
}
