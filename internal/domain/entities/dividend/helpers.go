package dividend

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
)

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

// Generates the future dividend payouts for a portfolio
func MakePayoutsFromDividendsAndLots(lots []lot.Lot, dividends []Dividend) []Payout {
	payouts := []Payout{}

	for i := range dividends {
		foundLot := false
		for j := range lots {
			if lots[j].Figi != dividends[i].Figi {
				continue
			}
			foundLot = true

			payout := Payout{
				Id:        uuid.New(),
				Figi:      lots[j].Figi,
				Ticker:    dividends[i].Security.Ticker,
				AccountId: lots[j].AccountId,
				Amount:    lots[j].Quantity * dividends[i].ExpectedDPS,
				Date:      dividends[i].PayoutDate,
			}
			payouts = append(payouts, payout)

		}
		if !foundLot {
			logger.Log("Failed to find the corresponding lot for dividend "+dividends[i].Id.String(), logger.ERROR)
		}
	}
	return payouts
}
