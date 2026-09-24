package bonds

import (
	"errors"
	"strings"

	"github.com/compoundinvest/invest-core/quote/entity"
	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
)

// Given a currency, returns all currency pairs present in the provided bond lots
// The provided lots must be populated with their corresponding bonds
func GetCurrencyPairs(cur2 string, lots []BondLot) []string {
	pairs := []string{}
	for i := range lots {
		if strings.EqualFold(lots[i].Bond.NominalCurrency, cur2) {
			continue
		}
		pairs = append(pairs, strings.Join([]string{lots[i].Bond.NominalCurrency, cur2}, "/"))
	}
	pairs = stringhelpers.RemoveDuplicatesFrom(pairs)
	return pairs
}

func MatchBondWithQuotes(bonds []Bond, quotes []entity.BondQuote) ([]Bond, []error) {
	var errorList []error
	for i := range bonds {
		foundQuote := false
		for j := range quotes {
			if quotes[j].GetTicker() != bonds[i].Ticker {
				continue
			}
			foundQuote = true
			bonds[i].QuoteInPercentage = quotes[j].GetQuoteAsPercentage()
		}
		if !foundQuote {
			errorList = append(errorList, errors.New("Failed to find a quote for bond lot "+bonds[i].Isin))
		}
	}
	return bonds, errorList
}
