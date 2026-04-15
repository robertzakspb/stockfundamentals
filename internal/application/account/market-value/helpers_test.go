package accountmvservice

import (
	"testing"

	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_ExtractMarketValueCurrencies(t *testing.T) {
	map1 := map[string]accountmvdomain.AccountMarketValue{
		"RUB": {},
		"USD": {},
	}
	map2 := map[string]accountmvdomain.AccountMarketValue{
		"RUB": {},
		"EUR": {},
		"USD": {},
	}

	currencies := ExtractMarketValueCurrencies(map1, map2)

	test.AssertEqual(t, 3, len(currencies))
	test.AssertEqual(t, "RUB", currencies[0])
	test.AssertEqual(t, "USD", currencies[1])
	test.AssertEqual(t, "EUR", currencies[2])
}
