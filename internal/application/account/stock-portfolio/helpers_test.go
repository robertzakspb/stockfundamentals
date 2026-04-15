package portfolio

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_GroupByNominalCurrency_Positive(t *testing.T) {
	lots := []lot.Lot{
		{
			Figi:     "lot1",
			Currency: "USD",
			Quantity: 10,
		},
		{
			Figi:     "lot1",
			Currency: "USD",
			Quantity: 5,
		},
		{
			Figi:     "lot1",
			Currency: "RUB",
			Quantity: 4,
		},
		{
			Figi:     "lot1",
			Currency: "RUB",
			Quantity: 3,
		},
		{
			Figi:     "lot1",
			Currency: "RUB",
			Quantity: 19,
		},
	}

	groupedLots := GroupByNominalCurrency(lots)

	test.AssertEqual(t, 10, groupedLots["USD"][0].Quantity)
	test.AssertEqual(t, 5, groupedLots["USD"][1].Quantity)
	test.AssertEqual(t, 4, groupedLots["RUB"][0].Quantity)
	test.AssertEqual(t, 3, groupedLots["RUB"][1].Quantity)
	test.AssertEqual(t, 19, groupedLots["RUB"][2].Quantity)
	test.AssertEqual(t, 2, len(groupedLots))
}

func Test_GroupByNominalCurrency_NoLots(t *testing.T) {
	lots := []lot.Lot{}

	groupedLots := GroupByNominalCurrency(lots)

	test.AssertEqual(t, 0, len(groupedLots))
}
