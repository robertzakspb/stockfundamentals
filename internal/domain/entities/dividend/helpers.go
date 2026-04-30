package dividend

import "github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"

func MatchDividendsWithStocks(divs []Dividend, securities []security.Stock) []Dividend {
	for i := range divs {
		for j := range securities {
			if divs[i].Figi == securities[j].Figi {
				divs[i].Security = securities[j]
			}
		}
	}

	return divs
}
