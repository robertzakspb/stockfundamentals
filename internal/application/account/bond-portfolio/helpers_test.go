package bondportfolio

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_MatchLotsWithBonds_ByIsin(t *testing.T) {
	lot := bonds.BondLot{
		Isin: "testIsin",
	}
	bond := bonds.Bond{
		Isin: "testIsin",
	}
	matchedLots := matchLotsWithBonds([]bonds.BondLot{lot}, []bonds.Bond{bond})

	test.AssertEqual(t, bond.Isin, matchedLots[0].Bond.Isin)
}

func Test_MatchLotsWithBonds_ByFigi(t *testing.T) {
	lot := bonds.BondLot{
		Figi: "testFigi",
	}
	bond := bonds.Bond{
		Figi: "testFigi",
	}
	matchedLots := matchLotsWithBonds([]bonds.BondLot{lot}, []bonds.Bond{bond})

	test.AssertEqual(t, bond.Figi, matchedLots[0].Bond.Figi)
}
