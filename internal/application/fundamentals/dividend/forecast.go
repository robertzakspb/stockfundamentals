package appdividend

import (
	"errors"
	"strconv"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/quotefetcher"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func SaveDividendForecast(forecast dividend.DividendForecast) error {

	figis, err := security_master.GetSecuritiesFilteredByFigi([]string{forecast.Stock.Figi})
	if err != nil {
		return err
	}
	if len(figis) == 0 {
		return errors.New("Failed to save the dividend forecast due to a missing corresponding security for figi " + forecast.Stock.Figi)
	}

	err = dbdividend.SaveDividendForecastToDb(mapDividendForecastToDbModel([]dividend.DividendForecast{forecast})[0])
	return err
}

func GetDividendForecasts() ([]dividend.DividendForecast, error) {
	dbForecasts, err := dbdividend.GetDividendForecasts()
	if err != nil {
		return []dividend.DividendForecast{}, err
	}
	forecasts := mapDividendForecastDbModelToDomain(dbForecasts)

	forecastsWithYields := populateYieldsForForecasts(forecasts)

	return forecastsWithYields, nil
}

func GetDivForecastsGroupedBySecurity() ([]dividend.SecurityDivForecast, error) {
	forecasts, err := GetDividendForecasts()
	if err != nil {
		return []dividend.SecurityDivForecast{}, err
	}

	securityDivForecasts := dividend.GroupForecastsBySecurity(forecasts)
	return securityDivForecasts, nil
}

func populateYieldsForForecasts(forecasts []dividend.DividendForecast) []dividend.DividendForecast {
	if len(forecasts) == 0 {
		return forecasts
	}

	securities := []entity.Security{}
	for _, forecast := range forecasts {
		security := entity.Security{
			Figi:   forecast.Stock.GetFigi(),
			MIC:    forecast.Stock.GetMic(),
			Ticker: forecast.Stock.GetTicker(),
			ISIN:   forecast.Stock.GetIsin(),
		}
		securities = append(securities, security)
	}
	quotes := quotefetcher.FetchQuotesFor(securities)
	if len(quotes) == 0 {
		logger.Log("Fetched 0 quotes for "+strconv.Itoa(len(forecasts))+" dividend forecasts", logger.ERROR)
	}

	for i := range forecasts {
		for _, quote := range quotes {
			if quote.Figi() != forecasts[i].Stock.Figi || quote.Quote() == 0 {
				continue
			}
			forecasts[i].Yield = forecasts[i].ExpectedDPS / quote.Quote()
		}
	}

	return forecasts
}
