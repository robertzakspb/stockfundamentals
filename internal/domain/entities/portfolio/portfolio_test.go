package stockportfolio

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_CollapseLotsIntoPositions_Negative_NothingToCollapse(t *testing.T) {
	lots := []lot.Lot{
		{
			Figi:         "figi1",
			Quantity:     10,
			PricePerUnit: 105,
		},
		{
			Figi:         "figi2",
			Quantity:     20,
			PricePerUnit: 130,
		},
	}

	collapsedPositions := CollapseLotsIntoPositionsUsingFigi(lots)

	test.AssertEqual(t, 2, len(collapsedPositions))
	test.AssertEqual(t, "figi1", collapsedPositions[0].Figi)
	test.AssertEqual(t, "figi2", collapsedPositions[1].Figi)
	test.AssertEqual(t, 10, collapsedPositions[0].Quantity)
	test.AssertEqual(t, 20, collapsedPositions[1].Quantity)
	test.AssertEqual(t, 105, collapsedPositions[0].PricePerUnit)
	test.AssertEqual(t, 130, collapsedPositions[1].PricePerUnit)
}

func Test_CollapseLotsIntoPositions_Positive_TwoCollapsedPositions(t *testing.T) {
	lots := []lot.Lot{
		{
			Figi:         "figi1",
			Quantity:     10,
			PricePerUnit: 105,
		},
		{
			Figi:         "figi1",
			Quantity:     20,
			PricePerUnit: 130,
		},
		{
			Figi:         "figi2",
			Quantity:     7,
			PricePerUnit: 97,
		},
		{
			Figi:         "figi2",
			Quantity:     27,
			PricePerUnit: 190,
		},
		{
			Figi:         "figi3",
			Quantity:     1,
			PricePerUnit: 5,
		},
	}

	collapsedPositions := CollapseLotsIntoPositionsUsingFigi(lots)

	test.AssertEqual(t, 3, len(collapsedPositions))
	test.AssertEqual(t, "figi1", collapsedPositions[0].Figi)
	test.AssertEqual(t, "figi2", collapsedPositions[1].Figi)
	test.AssertEqual(t, "figi3", collapsedPositions[2].Figi)
	test.AssertEqual(t, 30, collapsedPositions[0].Quantity)
	test.AssertEqual(t, (10.0*105.0+20.0*130.0)/30.0, collapsedPositions[0].PricePerUnit)
	test.AssertEqual(t, 34, collapsedPositions[1].Quantity)
	test.AssertEqual(t, (7.0*97.0+27.0*190.0)/34.0, collapsedPositions[1].PricePerUnit)
	test.AssertEqual(t, 1, collapsedPositions[2].Quantity)
	test.AssertEqual(t, 5, collapsedPositions[2].PricePerUnit)
}

func Test_GroupLotsByAccount_OneAccount(t *testing.T) {
	id := uuid.New()
	lots := []lot.Lot{
		{
			Figi:         "figi1",
			Quantity:     10,
			PricePerUnit: 105,
			AccountId:    id,
		},
		{
			Figi:         "figi1",
			Quantity:     20,
			PricePerUnit: 130,
			AccountId:    id,
		},
		{
			Figi:         "figi2",
			Quantity:     7,
			PricePerUnit: 97,
			AccountId:    id,
		},
		{
			Figi:         "figi2",
			Quantity:     27,
			PricePerUnit: 190,
			AccountId:    id,
		},
		{
			Figi:         "figi3",
			Quantity:     1,
			PricePerUnit: 5,
			AccountId:    id,
		},
	}

	accountLots := GroupLotsByAccount(lots)

	test.AssertEqual(t, 1, len(accountLots))
	test.AssertEqual(t, 10, accountLots[id][0].Quantity)
	test.AssertEqual(t, 20, accountLots[id][1].Quantity)
	test.AssertEqual(t, 7, accountLots[id][2].Quantity)
	test.AssertEqual(t, 27, accountLots[id][3].Quantity)
	test.AssertEqual(t, 1, accountLots[id][4].Quantity)
}

func Test_GroupLotsByAccount_ThreeAccounts(t *testing.T) {
	id1, id2, id3 := uuid.New(), uuid.New(), uuid.New()
	lots := []lot.Lot{
		{
			Figi:         "figi1",
			Quantity:     10,
			PricePerUnit: 105,
			AccountId:    id2,
		},
		{
			Figi:         "figi1",
			Quantity:     20,
			PricePerUnit: 130,
			AccountId:    id1,
		},
		{
			Figi:         "figi2",
			Quantity:     7,
			PricePerUnit: 97,
			AccountId:    id2,
		},
		{
			Figi:         "figi2",
			Quantity:     27,
			PricePerUnit: 190,
			AccountId:    id3,
		},
		{
			Figi:         "figi3",
			Quantity:     1,
			PricePerUnit: 5,
			AccountId:    id1,
		},
	}

	accountLots := GroupLotsByAccount(lots)

	test.AssertEqual(t, 3, len(accountLots))
	test.AssertEqual(t, 20, accountLots[id1][0].Quantity)
	test.AssertEqual(t, 1, accountLots[id1][1].Quantity)
	test.AssertEqual(t, 10, accountLots[id2][0].Quantity)
	test.AssertEqual(t, 7, accountLots[id2][1].Quantity)
	test.AssertEqual(t, 27, accountLots[id3][0].Quantity)
}
