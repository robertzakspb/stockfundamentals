package forexservice

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
	//"RUB": "", //FIXME!
}
