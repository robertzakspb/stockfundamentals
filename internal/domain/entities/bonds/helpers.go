package bonds

import (
	"strings"

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
