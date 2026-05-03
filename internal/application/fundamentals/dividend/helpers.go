package appdividend

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/google/uuid"
)

func matchDivForecastsWithPositions(forecasts []dividend.DividendForecast, positions []lot.Lot) ([]dividend.Payout, error) {
	payouts := []dividend.Payout{}

	for i := range forecasts {
		for j := range positions {
			if forecasts[i].Stock.Figi == positions[j].Figi {
				payout := dividend.Payout{
					Id:        uuid.New(),
					Figi:      positions[j].Figi,
					Ticker:    forecasts[i].Stock.Ticker,
					AccountId: positions[j].AccountId,
					Amount:    positions[j].Quantity * forecasts[i].ExpectedDPS,
					Date:      forecasts[i].ExpectedPayoutDate,
				}
				payouts = append(payouts, payout)
			}
		}
	}

	return payouts, nil
}
