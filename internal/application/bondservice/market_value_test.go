package bondservice

import (
	"testing"

	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_calculateRubMarketValue(t *testing.T) {
	bondList := []bonds.Bond{
		{
			Ticker:          "ticker1",
			NominalCurrency: "RUB",
			AccruedInterest: 12,
			NominalValue:    1000,
		},
		{
			Ticker:          "ticker2",
			NominalCurrency: "USD",
			AccruedInterest: 15,
			NominalValue:    1000,
		},
	}

	quotes := []tquoteservice.BondQuote{
		{QuoteAsPercentage: 95, Ticker: "ticker1"},
		{QuoteAsPercentage: 102, Ticker: "ticker2"},
		{QuoteAsPercentage: 34, Ticker: "ticker3"},
	}

	rates := []forexservice.ForexRate{
		{
			Currency1: "USD",
			Currency2: "RUB",
			Rate:      80,
		},
		{
			Currency1: "EUR",
			Currency2: "RUB",
			Rate:      90,
		},
	}

	bondList = calculateRubMarketValue(bondList, quotes, rates)

	test.AssertEqual(t, 2, len(bondList))
	test.AssertEqual(t, 962, bondList[0].MarketValueInRUB)
	test.AssertEqual(t, 82800, bondList[1].MarketValueInRUB)
}
