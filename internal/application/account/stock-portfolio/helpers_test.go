package portfolio

import (
	"testing"

	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
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

func Test_FindPortfolioByAccountId_Positive(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	portfolios := []stockportfolio.Portfolio{
		{
			Lots: []lot.Lot{
				{
					AccountId: id1,
				},
			},
		},
		{
			Lots: []lot.Lot{
				{
					AccountId: id2,
				},
			},
		},
	}

	portfolio, err := FindPortfolioByAccountId(id2, portfolios)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 1, len(portfolio.Lots))
	test.AssertEqual(t, id2, portfolio.Lots[0].AccountId)
}

func Test_FindPortfolioByAccountId_Negative(t *testing.T) {
	id1, id2, id3 := uuid.New(), uuid.New(), uuid.New()
	portfolios := []stockportfolio.Portfolio{
		{
			Lots: []lot.Lot{
				{
					AccountId: id1,
				},
			},
		},
		{
			Lots: []lot.Lot{
				{
					AccountId: id2,
				},
			},
		},
	}

	_, err := FindPortfolioByAccountId(id3, portfolios)

	test.AssertError(t, err)
}

func Test_GroupLotsByAccount(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	lots := []lot.Lot{
		{
			AccountId: id1,
			Quantity:  10,
		},
		{

			AccountId: id2,
			Quantity:  25,
		},
		{
			AccountId: id1,
			Quantity:  50,
		},
	}

	grouped := GroupLotsByAccount(lots)

	test.AssertEqual(t, 2, len(grouped))
	test.AssertEqual(t, 2, len(grouped[id1]))
	test.AssertEqual(t, 1, len(grouped[id2]))
	test.AssertEqual(t, id1, grouped[id1][0].AccountId)
	test.AssertEqual(t, 10, grouped[id1][0].Quantity)
	test.AssertEqual(t, id1, grouped[id1][1].AccountId)
	test.AssertEqual(t, 50, grouped[id1][1].Quantity)
	test.AssertEqual(t, id2, grouped[id2][0].AccountId)
	test.AssertEqual(t, 25, grouped[id2][0].Quantity)
}
