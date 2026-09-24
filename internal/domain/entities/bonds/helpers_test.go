package bonds

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_GetCurrencyPairs_Positive(t *testing.T) {
	lots := []BondLot{
		{
			Quantity: 5,
			Bond:     Bond{NominalCurrency: "RUB"},
		},
		{
			Quantity: 15,
			Bond:     Bond{NominalCurrency: "RUB"},
		},
		{
			Quantity: 10,
			Bond:     Bond{NominalCurrency: "USD"},
		},
		{
			Quantity: 12,
			Bond:     Bond{NominalCurrency: "EUR"},
		},
		{
			Quantity: 25,
			Bond:     Bond{NominalCurrency: "RUB"},
		},
	}

	pairs := GetCurrencyPairs("USD", lots)

	test.AssertEqual(t, 2, len(pairs))
	test.AssertEqual(t, "RUB/USD", pairs[0])
	test.AssertEqual(t, "EUR/USD", pairs[1])
}
