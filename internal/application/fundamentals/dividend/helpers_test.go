package appdividend

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_matchDivForecastsWithPositions(t *testing.T) {
	payoutDate := time.Now().AddDate(0, 0, 14)
	forecasts := []dividend.DividendForecast{
		{
			Id: uuid.New(),
			Stock: security.Stock{
				Figi:   "testFigi1",
				Ticker: "testTicker1",
			},
			ExpectedDPS:        25,
			Currency:           "USD",
			PaymentPeriod:      "H22026",
			Author:             "robert.z",
			Comment:            "testComment",
			Yield:              0.13,
			ExpectedPayoutDate: payoutDate,
		},
		{
			Id: uuid.New(),
			Stock: security.Stock{
				Figi:   "testFigi2",
				Ticker: "testTicker2",
			},
			ExpectedDPS:        10,
			Currency:           "USD",
			PaymentPeriod:      "H22026",
			Author:             "robert.z",
			Comment:            "testComment",
			Yield:              0.1,
			ExpectedPayoutDate: payoutDate,
		},
	}

	positions := []lot.Lot{
		{
			Quantity: 60,
			Figi:     "testFigi2",
		},
		{
			Quantity: 10,
			Figi:     "testFigi1",
		},
	}

	forecastedPayouts, err := matchDivForecastsWithPositions(forecasts, positions)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 2, len(forecastedPayouts))
	test.AssertEqual(t, "testFigi1", forecastedPayouts[0].Figi)
	test.AssertEqual(t, "testTicker1", forecastedPayouts[0].Ticker)
	test.AssertEqual(t, 250, forecastedPayouts[0].Amount)
	test.AssertEqual(t, payoutDate, forecastedPayouts[0].Date)
	test.AssertEqual(t, "testFigi2", forecastedPayouts[1].Figi)
	test.AssertEqual(t, "testTicker2", forecastedPayouts[1].Ticker)
	test.AssertEqual(t, 600, forecastedPayouts[1].Amount)
	test.AssertEqual(t, payoutDate, forecastedPayouts[1].Date)
}
