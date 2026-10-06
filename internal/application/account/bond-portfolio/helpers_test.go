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

func Test_GetLotTickers_Positive_MissingTicker(t *testing.T) {
	lots := []bonds.BondLot{
		{
			Bond: bonds.Bond{Ticker: "testTicker1"},
		},
		{
			Bond: bonds.Bond{Ticker: ""},
		},
		{
			Bond: bonds.Bond{Ticker: "testTicker3"},
		},
	}
	tickers := GetLotTickers(lots)

	test.AssertEqual(t, 2, len(tickers))
	test.AssertEqual(t, "testTicker1", tickers[0])
	test.AssertEqual(t, "testTicker3", tickers[1])
}

func Test_GetLotTickers_Positive(t *testing.T) {
	lots := []bonds.BondLot{
		{
			Bond: bonds.Bond{Ticker: "testTicker1"},
		},
		{
			Bond: bonds.Bond{Ticker: "testTicker2"},
		},
		{
			Bond: bonds.Bond{Ticker: "testTicker3"},
		},
	}
	tickers := GetLotTickers(lots)

	test.AssertEqual(t, 3, len(tickers))
	test.AssertEqual(t, "testTicker1", tickers[0])
	test.AssertEqual(t, "testTicker2", tickers[1])
	test.AssertEqual(t, "testTicker3", tickers[2])
}
