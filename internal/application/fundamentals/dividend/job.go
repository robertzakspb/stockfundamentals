package appdividend

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	securitydb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func FetchAndSaveAllDividends() error {
	dividends := fetchDividendsForAllStocks()

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		return errors.New("Unable to fetch dividends due to internal configuration issues")
	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return errors.New("Unable to fetch dividends due to internal database issues")
	}

	err = dbdividend.SaveDividendsToDB(dividends, db)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return errors.New(err.Error())
	}

	return nil
}

func fetchDividendsForAllStocks() []dividend.Dividend {
	stocks, err := securitydb.GetAllSecuritiesFromDB()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	allDividends := []dividend.Dividend{}

	rateLimit := time.Second
	throttle := time.Tick(rateLimit)

	for _, stock := range stocks {
		switch stock.GetCountry() {
		case "RU":
			dividends := fetchTinkoffDividendsFor(stock)
			allDividends = append(allDividends, dividends...)
			logger.Log(fmt.Sprintf("Fetched %d dividends for %s", len(dividends), stock.CompanyName), logger.INFORMATION)
		default:
			logger.Log("No data provider may provide dividends for "+stock.GetCompanyName(), logger.INFORMATION)
			continue
		}
		<-throttle
	}
	logger.Log("Completed the dividend fetching job", logger.INFORMATION)

	return allDividends
}
