package forexservice

import "strings"

// Returns currencies' coresponding symbols (e.g. USD -> $)
// If the symbol is not found, the provided argument is returned
func GetCurrencySymbol(curency string) string {
	symbol, found := currencyToSymbolMap[curency]
	if found {
		return symbol
	}

	return curency
}

var currencyToSymbolMap = map[string]string{
	"USD": "$",
	"RUB": "₽",
}

func FindRate(cur1, cur2 string, rates []ForexRate) (ForexRate, bool) {
	for _, rate := range rates {
		if string(rate.Currency1) == strings.ToUpper(cur1) && string(rate.Currency2) == strings.ToUpper(cur2) {
			return rate, true
		}
	}

	return ForexRate{}, false
}
