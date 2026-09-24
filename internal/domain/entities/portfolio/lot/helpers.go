package lot

import (
	"strings"

	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
)

func FindLotIndicesByFigi(lots []Lot, figi string) []int {
	filteredLotIndices := []int{}

	for i := range lots {
		if lots[i].Figi == figi {
			filteredLotIndices = append(filteredLotIndices, i)
		}
	}

	return filteredLotIndices
}

// Given a currency, returns all currency pairs present in the provided stock lots
func GetCurrencyPairs(cur2 string, lots []Lot) []string {
	pairs := []string{}
	for i := range lots {
		if strings.EqualFold(lots[i].Currency, cur2) {
			continue
		}
		pairs = append(pairs, strings.Join([]string{lots[i].Currency, cur2}, "/"))
	}
	pairs = stringhelpers.RemoveDuplicatesFrom(pairs)
	return pairs
}
