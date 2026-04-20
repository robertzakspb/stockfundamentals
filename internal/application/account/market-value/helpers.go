package accountmvservice

import (
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
)

func ExtractMarketValueCurrencies(mvMap1, mvMap2 map[string]accountmvdomain.AccountMarketValue) []string {
	currencies := []string{}

	for currency := range mvMap1 {
		currencies = append(currencies, currency)
	}
	for currency := range mvMap2 {
		currencies = append(currencies, currency)
	}

	currenciesSansDuplicates := stringhelpers.RemoveDuplicatesFrom(currencies)
	return currenciesSansDuplicates
}

func MarketValueCurrencyPairs(targetCurrency string, mvs []accountmvdomain.AccountMarketValue) []string {
	currencies := []string{}

	for _, mv := range mvs {
		currencies = append(currencies, mv.Currency+"/"+targetCurrency)
	}

	currenciesSansDuplicates := stringhelpers.RemoveDuplicatesFrom(currencies)
	return currenciesSansDuplicates
}
