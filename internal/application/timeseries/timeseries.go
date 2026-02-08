package timeseries

import (
	"errors"
	"strconv"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	tinkoffapi "github.com/compoundinvest/invest-core/quote/tinkoffmd"

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

	rateLimit := time.Second / 2 //So as not not overload the Tinkoff API
	throttle := time.Tick(rateLimit)
	for _, quote := range latestQuotes {
		if quote.Country != "RU" {
			continue
		}
		config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
		if err != nil {
			return errors.New("Unable to fetch dividends due to internal configuration issues")
		}

		tQuotes, _ := tinkoffapi.FetchAllHistoricalQuotesFor(quote.Figi, config, quote.Date, time.Now())
		logger.Log("Fetched "+strconv.Itoa(len(tQuotes))+" quotes for: "+quote.Figi, logger.INFORMATION)

		interfaceStructs := make([]entity.SimpleQuote, len(tQuotes))
		for i := range tQuotes {
			interfaceStructs[i] = tQuotes[i]
		}
		quotes = append(quotes, interfaceStructs...)

		if len(quotes) == 0 {
			continue
		}
		<-throttle
	}

	err = timeseries.SaveTimeSeriesToDB(quotes)
	if err != nil {
		logger.Log("Failed to save timeseries to DB due to: "+err.Error(), logger.ALERT)
		return errors.New("Failed to save timeseries to DB due to: " + err.Error())
	}

	return nil
}
