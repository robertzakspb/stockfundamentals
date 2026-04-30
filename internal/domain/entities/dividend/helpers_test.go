package dividend

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_matchDividendsWithStocks(t *testing.T) {
	dividends := []Dividend{
		{
			Figi: "figi1",
		},
		{
			Figi: "figi2",
		},
		{
			Figi: "figi3",
		},
	}

	stocks := []security.Stock {
		{
			Figi: "figi2",
		},
		{
			Figi: "figi1",
		},
		{
			Figi: "figi3",
		},
	}

	divsWithStocks := matchDividendsWithStocks(dividends, stocks)

	test.AssertEqual(t, 3, len(divsWithStocks))
	test.AssertEqual(t, "figi1", divsWithStocks[0].Security.Figi)
	test.AssertEqual(t, "figi2", divsWithStocks[1].Security.Figi)
	test.AssertEqual(t, "figi3", divsWithStocks[2].Security.Figi)


}
