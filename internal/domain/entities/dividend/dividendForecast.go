package dividend

import (
	"errors"
	"sort"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/google/uuid"
)

type DividendForecast struct {
	Id            uuid.UUID `json:"id"`
	Stock         security.Stock
	ExpectedDPS   float64 `json:"expectedDPS"`
	Currency      string  `json:"currency"`
	PaymentPeriod string  `json:"paymentPeriod"`
	Author        string  `json:"forecastAuthor"`
	Comment       string  `json:"comment"`
	Yield         float64 `json:"yield"`
}

type SecurityDivForecast struct {
	Figi      string
	Forecasts []DividendForecast
}

func (f SecurityDivForecast) CumulativeReturn() float64 {
	totalValue := 1.0
	for _, forecast := range f.Forecasts {
		totalValue *= (1 + forecast.Yield)
	}
	totalReturn := totalValue - 1
	return totalReturn
}

func getForecastIndex(forecasts []SecurityDivForecast, figi string) (int, error) {
	for i, forecast := range forecasts {
		if forecast.Figi == figi {
			return i, nil
		}
	}
	return -1, errors.New("No security with figi " + figi + " was found in the div forecasts")
}

func GroupForecastsBySecurity(forecasts []DividendForecast) []SecurityDivForecast {
	secDivForecasts := []SecurityDivForecast{}
	for i := range forecasts {
		securityIndex, err := getForecastIndex(secDivForecasts, forecasts[i].Stock.Figi)
		if err != nil {
			secDivForecasts = append(secDivForecasts, SecurityDivForecast{
				Figi:      forecasts[i].Stock.Figi,
				Forecasts: []DividendForecast{forecasts[i]},
			})
			continue
		}
		secDivForecasts[securityIndex].Forecasts = append(secDivForecasts[securityIndex].Forecasts, forecasts[i])
	}
	sort.Slice(secDivForecasts, func(i, j int) bool {
		return secDivForecasts[i].CumulativeReturn() > secDivForecasts[j].CumulativeReturn()
	})

	return secDivForecasts
}

// func (f SecurityDivForecasts) AnnualizedReturn() float64 {
// 	totalReturn := f.TotalReturn()
// 	annualizedReturn := compoundinterest.CalcAnnualizedReturn(totalReturn, )
// }
