package forex

import "fmt"

func GetExchangeRateForPair(currency1, currency2 string, dp ForexDataProvider) (float64, error) {
	if currency1 == currency2 {
		return 1, nil
	}

	usdToCur1, err := dp.GetExchangeRateUsdTo(currency1)
	usdToCur2, err := dp.GetExchangeRateUsdTo(currency2)
	if err != nil || usdToCur2 == 0 {
		return 0, fmt.Errorf("%s", "failed to find the exchange rate between "+string(USD)+"and "+currency2)
	}

	return usdToCur1 / usdToCur2 , nil
}

func ConvertPriceToDifferentCurrency(price float64, priceCur, targetCur string, dp ForexDataProvider) (float64, error) {
	exRate, err := GetExchangeRateForPair( targetCur, priceCur, dp)
	if err != nil {
		return 0, err
	}

	return price / exRate, nil
}

var exchangeRateUsdTo = map[Currency]float64{
	RUB: 80,
	EUR: 0.86,
	RSD: 99.54,
}

type Currency string

const (
	USD Currency = "USD"
	RUB Currency = "RUB"
	EUR Currency = "EUR"
	RSD Currency = "RSD"
)

var currencyName = map[string]Currency{
	"USD": USD,
	"RUB": RUB,
	"EUR": EUR,
	"RSD": RSD,
}
