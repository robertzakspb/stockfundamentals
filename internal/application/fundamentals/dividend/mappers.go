package appdividend

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
)

func mapDividendForecastToDbModel(forecasts []dividend.DividendForecast) []dbdividend.DividendForecastDb {
	dbModels := []dbdividend.DividendForecastDb{}

	for _, forecast := range forecasts {
		dbModel := dbdividend.DividendForecastDb{
			Figi:          forecast.Figi,
			ExpectedDPS:   forecast.ExpectedDPS,
			Currency:      forecast.Currency,
			PaymentPeriod: forecast.PaymentPeriod,
			Author:        forecast.Author,
			Comment:       forecast.Author,
		}

		dbModels = append(dbModels, dbModel)
	}
	return dbModels
}
