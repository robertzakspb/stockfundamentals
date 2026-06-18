package forexservice

import (
	"strings"
	"time"
)

type ForexRate struct {
	Currency1 Currency  `sql:"currency_1" json:"currency1"`
	Currency2 Currency  `sql:"currency_2" json:"currency2"`
	Rate      float64   `sql:"rate" json:"rate"`
	Date      time.Time `sql:"date" json:"date"`
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
