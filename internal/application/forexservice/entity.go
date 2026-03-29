package forexservice

import (
	"strings"
	"time"
)

type ForexRate struct {
	Currency1 Currency
	Currency2 Currency
	Rate      float64
	Date      time.Time
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

func IsSupportedCurrency(currency string) bool {
	_, found := currencyName[strings.ToUpper(currency)]
	return found
}
