package forexservice

func StartFxRateImportJob() {
	var requiredCurrencyPairs = map[Currency]Currency{
		USD: RUB,
		EUR: RUB,
	}
	for cur1, cur2 := range requiredCurrencyPairs {
		FetchAndSaveCurrencyPairQuotes(string(cur1), string(cur2))
	}
}
