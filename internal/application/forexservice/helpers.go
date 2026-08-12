package forexservice

import (
	"errors"
	"slices"
	"strings"
)

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

// Takes currency pairs and generates a list of all currency 1s and currency 2s that will be used as filtes in the DB query
func generateCurrency1AndCurrency2Slices(currencyPairs []string) (cur1s []string, cur2s []string) {
	for i, pair := range currencyPairs {
		split := strings.Split(pair, "/")
		if i == 0 { //We only need to add USD once to prevent the repetitive "WHERE cur1 = USD AND cur1 = USD etc."
			cur1s = append(cur1s, "USD")
		}
		//USD-to-X rates can be fetched directly from the DB
		if split[0] == "USD" {
			cur2 := strings.ToUpper(split[1])
			if !slices.Contains(cur2s, cur2) {
				cur2s = append(cur2s, cur2)
			}
		} else {
			// The app currently stores only the exchange rate of USD to other currencies. Hence, if a user wants to get
			// the exchange rate between two non-USD currencies, this function can calculate it through cross-rates
			cur1, cur2 := strings.ToUpper(split[0]), strings.ToUpper(split[1])
			if !slices.Contains(cur2s, cur1) {
				cur2s = append(cur2s, cur1) //Looking for the USD-to-first-currency rate
			}
			if !slices.Contains(cur2s, cur2) {
				cur2s = append(cur2s, cur2) //Looking for the USD-to-second-currency rate
			}
		}
	}

	return cur1s, cur2s
}

// Calculates the cross rate between two currencies given their USD rates
func calculateCrossRateViaUsdRates(usdRate1, usdRate2 ForexRate) (ForexRate, error) {
	if usdRate1.Rate == 0 {
		return ForexRate{}, errors.New("The exchange rate is 0 for the first rate")
	}

	crossRateValue := usdRate2.Rate / usdRate1.Rate
	crossRate := ForexRate{
		Currency1: usdRate1.Currency2,
		Currency2: usdRate2.Currency2,
		Rate:      crossRateValue,
		Date:      usdRate1.Date,
	}

	return crossRate, nil
}

// If, for instance, the EUR/RSD rate is required, the app will fetch two rates: USD/RSD and USD/EUR.
// These two rates must be collapsed into the cross rate (EUR/RSD) when returning the values
func collapseRatesIntoTargetCrossRates(currencyPairs []string, rates []ForexRate) ([]ForexRate, error) {
	cleanRates := []ForexRate{}

	//Adding the USD/X rates first
	for i := range rates {
		if rates[i].Currency1 == "USD" && slices.Contains(currencyPairs, strings.Join([]string{string(rates[i].Currency1), string(rates[i].Currency2)}, "/")) {
			cleanRates = append(cleanRates, rates[i])
		}
	}

	//Proceeding to collapse the rates into cross rates
	crossRatePairs := []string{}
	for _, pair := range currencyPairs {
		if pair[:3] != "USD" {
			crossRatePairs = append(crossRatePairs, pair)
		}
	}

	for _, crossRatePair := range crossRatePairs {
		split := strings.Split(crossRatePair, "/")
		rate1, found := FindRate("USD", split[0], rates)
		if !found {
			return cleanRates, errors.New("Failed to find the USD/" + split[0] + " rate in the list")
		}
		rate2, found := FindRate("USD", split[1], rates)
		if !found {
			return cleanRates, errors.New("Failed to find the USD/" + split[1] + " rate in the list")
		}
		crossRate, err := calculateCrossRateViaUsdRates(rate1, rate2)
		if err != nil {
			return cleanRates, err
		}
		cleanRates = append(cleanRates, crossRate)
	}

	return cleanRates, nil
}
