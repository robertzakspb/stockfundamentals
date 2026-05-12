package timeseries

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	tinkoffapi "github.com/compoundinvest/invest-core/quote/tinkoffmd"

	tthrottler "github.com/compoundinvest/stockfundamentals/internal/application/tinkoff-throttler"
	timeseriesdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
)

func FetchAndSaveHistoricalQuotes() error {
	latestQuotes, err := timeseriesdb.GetLatestQuotesForAllSecurities()
	if err != nil {
		return err
	}
	quotes := []entity.SimpleQuote{}

	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		return errors.New("Unable to fetch quotes due to internal configuration issues")
	}
	client, err := tinkoff.NewClient(context.TODO(), config, nil)
	service := client.NewMarketDataServiceClient()

	for i, quote := range latestQuotes {
		<-tthrottler.InstrumentServiceThrottle

		if quote.Country != "RU" {
			continue
		}

		tQuotes, err := tinkoffapi.FetchAllHistoricalQuotesFor(service, quote.Figi, quote.Date, time.Now())
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
		}
		logger.Log(strconv.Itoa(i+1)+" out of "+strconv.Itoa(len(latestQuotes))+". Fetched "+strconv.Itoa(len(tQuotes))+" quotes for: "+quote.Figi, logger.INFORMATION)

		interfaceStructs := make([]entity.SimpleQuote, len(tQuotes))
		for i := range tQuotes {
			interfaceStructs[i] = tQuotes[i]
		}
		quotes = interfaceStructs

		if len(quotes) == 0 {
			continue
		}
		go timeseriesdb.SaveTimeSeriesToDB(&quotes)
	}

	logger.Log("The time series job has successfully completed", logger.INFORMATION)

	return nil
}
