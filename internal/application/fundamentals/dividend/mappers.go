package appdividend

import (
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func mapDividendForecastToDbModel(forecasts []dividend.DividendForecast) []dbdividend.DividendForecastDb {
	dbModels := []dbdividend.DividendForecastDb{}

	for _, forecast := range forecasts {
		dbModel := dbdividend.DividendForecastDb{
			Id:            forecast.Id,
			Figi:          forecast.Stock.Figi,
			ExpectedDPS:   forecast.ExpectedDPS,
			Currency:      forecast.Currency,
			PaymentPeriod: forecast.PaymentPeriod,
			Author:        forecast.Author,
			Comment:       forecast.Comment,
		}

		dbModels = append(dbModels, dbModel)
	}
	return dbModels
}

func mapDividendForecastDbModelToDomain(dbModels []dbdividend.DividendForecastDb) []dividend.DividendForecast {
	forecasts := []dividend.DividendForecast{}

	figis := []string{}
	for _, forecast := range forecasts {
		figis = append(figis, forecast.Stock.Figi)
	}
	securities, err := security_master.GetSecuritiesFilteredByFigi(figis)
	if err != nil {
		logger.Log("Failed to fetch securities with the target Figis", logger.ERROR)
	}

	for _, forecast := range dbModels {
		var targetStock security.Stock
		for _, stock := range securities {
			if stock.Figi != forecast.Figi {
				continue
			}
			targetStock = stock
		}
		if targetStock.Figi == "" {
			targetStock.Figi = forecast.Figi
			logger.Log("Failed to find a security with figi "+forecast.Figi, logger.ERROR)
		}

		domainStruct := dividend.DividendForecast{
			Id:            forecast.Id,
			Stock:         targetStock,
			ExpectedDPS:   forecast.ExpectedDPS,
			Currency:      forecast.Currency,
			PaymentPeriod: forecast.PaymentPeriod,
			Author:        forecast.Author,
			Comment:       forecast.Author,
		}

		forecasts = append(forecasts, domainStruct)
	}
	return forecasts
}
