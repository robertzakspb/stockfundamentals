package timeseries

import (
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	tinkoffapi "github.com/compoundinvest/invest-core/quote/tinkoffmd"
	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/security"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	tinkoff "github.com/russianinvestments/invest-api-go-sdk/investgo"
)

// type HistoricalReturn struct {
// 	Ticker      string
// 	FirstQuote  tinkoffapi.TinkoffQuote
// 	LastQuote   tinkoffapi.TinkoffQuote
// 	TotalReturn float64
// }

func FetchAndSaveHistoricalQuotes() {
	stocks, _ := security.FetchSecuritiesFromDB()
	quotes := []entity.SimpleQuote{}

	rateLimit := time.Second / 2
	throttle := time.Tick(rateLimit)
	for _, stock := range stocks { 
		if stock.Country != "RU" {
			continue
		}
		config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
		if err != nil {
			panic("Unable to initialize the Tinkoff API config file")
		}

		tQuotes, _ := tinkoffapi.FetchAllHistoricalQuotesFor(entity.Security{Figi: stock.Figi, ISIN: stock.Isin}, config)

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

	err := saveTimeSeriesToDB(quotes)
	if err != nil {
		logger.Log("Failed to fetch timeseries via Tinkoff API due to: " + err.Error(), logger.ALERT)
	}
}
