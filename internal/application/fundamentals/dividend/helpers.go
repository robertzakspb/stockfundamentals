package appdividend

import (
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
)

func PopulateDividendStocks(dividends []dividend.Dividend) ([]dividend.Dividend, error) {
	figis := extractFigisFromDividends(dividends)
	stocks, err := security_master.GetSecuritiesFilteredByFigi(figis)
	if err != nil {
		return dividends, err
	}

	for i := range dividends {
		foundStock := false
		for j := range stocks {
			if stocks[j].Figi != dividends[i].Figi {
				continue
			}
			foundStock = true
			dividends[i].Security = stocks[j]
		}
		if !foundStock {
			logger.Log("Failed to find the stock for dividend "+dividends[i].Figi, logger.ERROR)
		}
	}

	return dividends, nil
}

func extractFigisFromDividends(dividends []dividend.Dividend) []string {
	figis := []string{}
	for i := range dividends {
		if dividends[i].Figi == "" {
			continue
		}
		figis = append(figis, dividends[i].Figi)
	}

	return figis
}

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

// func matchDividendPayoutsWithPositions(dividends []dividend.Dividend, positions []lot.Lot) ([]dividend.Payout, error) {
// 	payouts := []dividend.Payout{}

// 	for i := range dividends {
// 		for j := range positions {
// 			if dividends[i].Security.Figi == positions[j].Figi {
// 				payout := dividend.Payout{
// 					Id:        uuid.New(),
// 					Figi:      positions[j].Figi,
// 					Ticker:    dividends[i].Security.Ticker,
// 					AccountId: positions[j].AccountId,
// 					Amount:    positions[j].Quantity * dividends[i].ExpectedDPS,
// 					Date:      dividends[i].PayoutDate,
// 				}
// 				payouts = append(payouts, payout)
// 			}
// 		}
// 	}

// 	return payouts, nil
// }
