package bondservice

import (
	"testing"
	"time"

	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_calculateRubMarketValue(t *testing.T) {
	bondList := []bonds.Bond{
		{
			Figi:            "figi1",
			NominalCurrency: "RUB",
			AccruedInterest: 12,
			NominalValue:    1000,
		},
		{
			Figi:            "figi2",
			NominalCurrency: "USD",
			AccruedInterest: 15,
			NominalValue:    1000,
		},
	}

	quotes := []tquoteservice.TQuote{
		tquoteservice.New(95, "figi1", "ticker1", time.Now()),
		tquoteservice.New(102, "figi2", "ticker2", time.Now()),
		tquoteservice.New(34, "figi3", "ticker3", time.Now()),
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
