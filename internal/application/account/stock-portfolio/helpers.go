package portfolio

import "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"

func GroupByNominalCurrency(lots []lot.Lot) map[string][]lot.Lot {
	groupedLots := map[string][]lot.Lot{}

	for i := range lots {
		groupedLots[lots[i].Currency] = append(groupedLots[lots[i].Currency], lots[i])
	}

	return groupedLots
}
