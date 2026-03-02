package appdividend

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
)

func SaveDividendForecast(forecast dividend.DividendForecast) error {
	err := dbdividend.SaveDividendForecastToDb(mapDividendForecastToDbModel([]dividend.DividendForecast{forecast})[0])
	return err
}