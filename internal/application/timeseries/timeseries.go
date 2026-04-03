package timeseries

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	tinkoffapi "github.com/compoundinvest/invest-core/quote/tinkoffmd"

	tthrottler "github.com/compoundinvest/stockfundamentals/internal/application/tinkoff-throttler"
	timeseries "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
)

func FetchAndSaveHistoricalQuotes() error {
	latestQuotes, err := timeseries.GetLatestQuotesForAllSecurities()
	if err != nil {
		return err
	}
	quotes := []entity.SimpleQuote{}

	for _, quote := range latestQuotes {
		if quote.Country != "RU" {
			continue
		}
		config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
		if err != nil {
			return errors.New("Unable to fetch dividends due to internal configuration issues")
		}

		client, err := tinkoff.NewClient(context.TODO(), config, nil)
		service := client.NewMarketDataServiceClient()

		tQuotes, err := tinkoffapi.FetchAllHistoricalQuotesFor(service, quote.Figi, quote.Date, time.Now())
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
		}
		logger.Log("Fetched "+strconv.Itoa(len(tQuotes))+" quotes for: "+quote.Figi, logger.INFORMATION)

		interfaceStructs := make([]entity.SimpleQuote, len(tQuotes))
		for i := range tQuotes {
			interfaceStructs[i] = tQuotes[i]
		}
		quotes = append(quotes, interfaceStructs...)

		if len(quotes) == 0 {
			continue
		}
		err = timeseries.SaveTimeSeriesToDB(quotes)
		if err != nil {
			logger.Log("Failed to save timeseries for "+quote.Figi+" to DB due to: "+err.Error(), logger.ALERT)
			continue
		}

		<-tthrottler.MarketDataServiceThrottle
	}

	return nil
}
