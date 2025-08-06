package forex

import (
	"fmt"
	"strings"
)

func GetExchangeRateUsdTo(currency string) (float64, error) {
	if currency == "USD" {
		return 1, nil
	}

	cur, found := currencyName[currency]
	if !found {
		return 0, fmt.Errorf("%s", "Exchange rates are currently unavailable for "+currency)
	}

	//TODO: Refactor this to dynamically update currencies on a daily basis
	rate, found := exchangeRateUsdTo[cur]
	if !found {
		return 0, fmt.Errorf("%s", "failed to find the exchange rate between "+string(USD)+"and "+currency)
	}

	return rate, nil
}

func GetExchangeRateForPair(currency1, currency2 string) (float64, error) {
	if currency1 == currency2 {
		return 1, nil
	}

	usdToCur1, err := GetExchangeRateUsdTo(currency1)
	usdToCur2, err := GetExchangeRateUsdTo(currency2)
	if err != nil || usdToCur2 == 0 {
		return 0, fmt.Errorf("%s", "failed to find the exchange rate between "+string(USD)+"and "+currency2)
	}

	return usdToCur2 / usdToCur1, nil
}

func ConvertPriceToDifferentCurrency(price float64, priceCur string, targetCur string) (float64, error) {
	exRate, err := GetExchangeRateForPair(targetCur, priceCur)
	if err != nil {
		return 0, err
	}

	return price / exRate, nil
}

func IsSupportedCurrency(cur string) bool {
	_, found := currencyName[strings.ToUpper(cur)]
	return found
}

var exchangeRateUsdTo = map[Currency]float64{
	RUB: 79,
	EUR: 0.85,
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
