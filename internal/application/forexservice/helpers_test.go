package forexservice

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_generateCurrency1AndCurrency2Slices_OnlyUSD(t *testing.T) {
	currencyPairs := []string{"USD/RUB", "USD/EUR", "USD/RSD"}

	cur1s, cur2s := generateCurrency1AndCurrency2Slices(currencyPairs)

	test.AssertEqual(t, 1, len(cur1s))
	test.AssertEqual(t, 3, len(cur2s))

	test.AssertEqual(t, "USD", cur1s[0])
	test.AssertEqual(t, "RUB", cur2s[0])
	test.AssertEqual(t, "EUR", cur2s[1])
	test.AssertEqual(t, "RSD", cur2s[2])
}

func Test_generateCurrency1AndCurrency2Slices_OnlyNonUSD(t *testing.T) {
	currencyPairs := []string{"CHF/RUB", "RSD/RUB", "EUR/RUB"}

	cur1s, cur2s := generateCurrency1AndCurrency2Slices(currencyPairs)

	test.AssertEqual(t, 1, len(cur1s))
	test.AssertEqual(t, 4, len(cur2s))

	test.AssertEqual(t, "USD", cur1s[0])
	test.AssertEqual(t, "CHF", cur2s[0])
	test.AssertEqual(t, "RUB", cur2s[1])
	test.AssertEqual(t, "RSD", cur2s[2])
	test.AssertEqual(t, "EUR", cur2s[3])
}

func Test_generateCurrency1AndCurrency2Slices_UsdAndNonUSD(t *testing.T) {
	currencyPairs := []string{"CHF/RUB", "RSD/RUB", "EUR/RUB", "USD/CNY", "USD/EUR"}

	cur1s, cur2s := generateCurrency1AndCurrency2Slices(currencyPairs)

	test.AssertEqual(t, 1, len(cur1s))
	test.AssertEqual(t, 5, len(cur2s))

	test.AssertEqual(t, "USD", cur1s[0])
	test.AssertEqual(t, "CHF", cur2s[0])
	test.AssertEqual(t, "RUB", cur2s[1])
	test.AssertEqual(t, "RSD", cur2s[2])
	test.AssertEqual(t, "EUR", cur2s[3])
	test.AssertEqual(t, "CNY", cur2s[4])
}

func Test_calculateRateViaCrossRates(t *testing.T) {
	date := time.Now()
	rate1 := ForexRate{
		Currency1: "USD",
		Currency2: "RUB",
		Rate:      80,
		Date:      date,
	}
	rate2 := ForexRate{
		Currency1: "USD",
		Currency2: "RSD",
		Rate:      100,
		Date:      date,
	}

	crossRate, err := calculateCrossRateViaUsdRates(rate1, rate2)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "RUB", crossRate.Currency1)
	test.AssertEqual(t, "RSD", crossRate.Currency2)
	test.AssertEqual(t, 100.0/80.0, crossRate.Rate)
	test.AssertEqual(t, date, crossRate.Date)
}

func Test_collapseRatesIntoTargetCrossRates_Positive(t *testing.T) {
	currencyPairs := []string{"USD/RUB", "USD/RSD", "RUB/RSD", "EUR/RUB"}
	rates := []ForexRate{
		{
			Currency1: "USD",
			Currency2: "RUB",
			Rate:      80,
		},
		{
			Currency1: "USD",
			Currency2: "RSD",
			Rate:      100,
		},
		{
			Currency1: "USD",
			Currency2: "EUR",
			Rate:      0.9,
		},
	}

	cleanRates, err := collapseRatesIntoTargetCrossRates(currencyPairs, rates)

	test.AssertNoError(t, err)

	test.AssertEqual(t, 4, len(cleanRates))
	test.AssertEqual(t, "USD", cleanRates[0].Currency1)
	test.AssertEqual(t, "RUB", cleanRates[0].Currency2)
	test.AssertEqual(t, 80, cleanRates[0].Rate)

	test.AssertEqual(t, "USD", cleanRates[1].Currency1)
	test.AssertEqual(t, "RSD", cleanRates[1].Currency2)
	test.AssertEqual(t, 100, cleanRates[1].Rate)

	test.AssertEqual(t, "RUB", cleanRates[2].Currency1)
	test.AssertEqual(t, "RSD", cleanRates[2].Currency2)
	test.AssertEqual(t, 100.0/80.0, cleanRates[2].Rate)

	test.AssertEqual(t, "EUR", cleanRates[3].Currency1)
	test.AssertEqual(t, "RUB", cleanRates[3].Currency2)
	test.AssertEqual(t, 80.0/0.9, cleanRates[3].Rate)

}
