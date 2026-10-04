package dividend

import (
	"math"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
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
func MakePayoutsFromDividendsAndLots(lots []lot.Lot, dividends []Dividend, transactions []transaction.Transaction) []Payout {
	payouts := []Payout{}

	for i := range dividends {
		foundLot := false
		for j := range lots {
			if lots[j].Figi != dividends[i].Figi {
				continue
			}
			foundLot = true

			hasBeenPaidOut := false
			for k := range transactions {
				if transactions[k].Figi != dividends[i].Figi {
					continue
				}
				differenceIsLessThanOnePercent := math.Abs((transactions[k].PricePerUnit-lots[j].Quantity*dividends[i].ActualDPS))/transactions[k].PricePerUnit < 0.01
				if differenceIsLessThanOnePercent {
					hasBeenPaidOut = true
				}
			}

			payout := Payout{
				Id:             uuid.New(),
				Figi:           lots[j].Figi,
				Ticker:         dividends[i].Security.Ticker,
				AccountId:      lots[j].AccountId,
				Amount:         lots[j].Quantity * dividends[i].ActualDPS,
				Date:           dividends[i].PayoutDate,
				Dividend:       dividends[i],
				HasBeenPaidOut: hasBeenPaidOut,
			}
			payouts = append(payouts, payout)

		}
		if !foundLot {
			logger.Log("Failed to find the corresponding lot for dividend "+dividends[i].Id.String(), logger.ERROR)
		}
	}
	return payouts
}
