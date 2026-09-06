package currency

import "errors"

var currency_divisor = map[string]int{
	"RUB": 100,
	"USD": 100,
	"EUR": 100,
	"RSD": 100,
}

func GetCurrencyDivisor(currency string) (int, error) {
	divisor, found := currency_divisor[currency]
	if !found {
		return -1, errors.New("Failed to find the divisor for currency " + currency)
	}

	return divisor, nil
}
