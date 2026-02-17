package forex

import (
	"fmt"
	"strings"
)

type ForexDataProvider interface {
	GetExchangeRateUsdTo(currency string) (float64, error)
	IsSupportedCurrency(cur string) bool
}

// Implementation of the Forex Data Provider interface
type ForexDP struct {
}

func (f ForexDP) GetExchangeRateUsdTo(currency string) (float64, error) {
	if strings.ToLower(currency) == "usd" {
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

func (f ForexDP) IsSupportedCurrency(cur string) bool {
	_, found := currencyName[strings.ToUpper(cur)]
	return found
}
