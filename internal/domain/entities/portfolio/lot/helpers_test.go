package lot

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_FindLotsByFigi(t *testing.T) {
	targetFigi := "figi1"
	lots := []Lot{
		{
			Figi: targetFigi,
		},
		{
			Figi: "figi5",
		},
		{
			Figi: targetFigi,
		},
		{
			Figi: "figi2",
		},
		{
			Figi: "figi3",
		},
		{
			Figi: targetFigi,
		},
	}

	filteredLots := FindLotIndicesByFigi(lots, targetFigi)

	test.AssertEqual(t, 3, len(filteredLots))
	test.AssertEqual(t, targetFigi, lots[filteredLots[0]].Figi)
	test.AssertEqual(t, targetFigi, lots[filteredLots[1]].Figi)
	test.AssertEqual(t, targetFigi, lots[filteredLots[2]].Figi)
}

func Test_GetCurrencyPairs_Positive(t *testing.T) {
	lots := []Lot{
		{
			Quantity: 5,
			Currency: "RUB",
		},
		{
			Quantity: 15,
			Currency: "RUB",
		},
		{
			Quantity: 10,
			Currency: "USD",
		},
		{
			Quantity: 12,
			Currency: "EUR",
		},
		{
			Quantity: 25,
			Currency: "RUB",
		},
	}

	pairs := GetCurrencyPairs("USD", lots)

	test.AssertEqual(t, 2, len(pairs))
	test.AssertEqual(t, "RUB/USD", pairs[0])
	test.AssertEqual(t, "EUR/USD", pairs[1])
}