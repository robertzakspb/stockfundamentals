package appdividend

import (
	"context"
	"fmt"

	tthrottler "github.com/compoundinvest/stockfundamentals/internal/application/tinkoff-throttler"

	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	securitydb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func FetchAndSaveAllDividends() error {
	stocks, err := securitydb.GetAllSecuritiesFromDB()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	client, err := tinkoff.NewClient(context.TODO(), config, nil)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}

	securityService := client.NewInstrumentsServiceClient()

	for _, stock := range stocks {
		switch stock.GetCountry() {
		case "RU":
			dividends := fetchTinkoffDividendsFor(securityService, stock)
			logger.Log(fmt.Sprintf("Fetched %d dividends for %s", len(dividends), stock.CompanyName), logger.INFORMATION)
			go dbdividend.SaveDividendsToDB(&dividends)

		default:
			logger.Log("No data provider may provide dividends for "+stock.GetCompanyName(), logger.INFORMATION)
			continue
		}
		<-tthrottler.InstrumentServiceThrottle
	}
	logger.Log("Completed the dividend fetching job", logger.INFORMATION)

	return nil
}
