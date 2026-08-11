package forexservice

import "strings"

func ImportForexRatesJob() {
	var requiredCurrencyPairs = []string{"USD/RUB", "EUR/RUB", "USD/RSD"}

	for _, currencyPair := range requiredCurrencyPairs {
		split := strings.Split(currencyPair, "/")
		cur1 := split[0]
		cur2 := split[1]
		FetchAndSaveCurrencyPairQuotes(string(cur1), string(cur2))
	}
}
