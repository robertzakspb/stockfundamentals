package accountmvservice

import (
	"slices"
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

	slices.SortFunc(currencies, func(cur1, cur2 string) int {
		if cur1 < cur2 {
			return -1
		} else {
			return 1
		}
	})

	test.AssertEqual(t, 3, len(currencies))
	test.AssertEqual(t, "EUR", currencies[0])
	test.AssertEqual(t, "RUB", currencies[1])
	test.AssertEqual(t, "USD", currencies[2])
}
