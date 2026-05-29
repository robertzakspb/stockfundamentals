package portfolio

import (
	"errors"

	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/google/uuid"
)

func GroupByNominalCurrency(lots []lot.Lot) map[string][]lot.Lot {
	groupedLots := map[string][]lot.Lot{}

	for i := range lots {
		groupedLots[lots[i].Currency] = append(groupedLots[lots[i].Currency], lots[i])
	}

	return groupedLots
}

func FindPortfolioByAccountId(id uuid.UUID, portfolios []stockportfolio.Portfolio) (stockportfolio.Portfolio, error) {
	for i := range portfolios {
		if len(portfolios[i].Lots) == 0 {
			continue
		}
		if portfolios[i].Lots[0].AccountId == id {
			return portfolios[i], nil
		}
	}
	return stockportfolio.Portfolio{}, errors.New("Failed to find the portfolio for account " + id.String())
}

func GroupLotsByAccount(lots []lot.Lot) map[uuid.UUID][]lot.Lot {
	groupedLots := map[uuid.UUID][]lot.Lot{}

	for i := range lots {
		groupedLots[lots[i].AccountId] = append(groupedLots[lots[i].AccountId], lots[i])
	}

	return groupedLots
}
