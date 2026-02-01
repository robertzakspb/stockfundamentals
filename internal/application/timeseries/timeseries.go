package timeseries

import (
	"net/http"
	"strconv"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	tinkoffapi "github.com/compoundinvest/invest-core/quote/tinkoffmd"
	"github.com/gin-gonic/gin"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
	securitydb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	tinkoff "github.com/russianinvestments/invest-api-go-sdk/investgo"
)

func FetchAndSaveHistoricalQuotes(c *gin.Context) {
	stocks, _ := securitydb.GetAllSecuritiesFromDB()
	quotes := []entity.SimpleQuote{}

	rateLimit := time.Second / 2 //So as not not overload the Tinkoff API
	throttle := time.Tick(rateLimit)
	for _, stock := range stocks {
		if stock.Country != "RU" {
			continue
		}
		config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Unable to fetch dividends due to internal configuration issues")
		}

		tQuotes, _ := tinkoffapi.FetchAllHistoricalQuotesFor(entity.Security{Figi: stock.Figi, ISIN: stock.Isin}, config)
		logger.Log("Fetched "+strconv.Itoa(len(tQuotes)) + " quotes for: "+stock.Ticker, logger.INFORMATION)

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

	err := timeseries.SaveTimeSeriesToDB(quotes)
	if err != nil {
		logger.Log("Failed to fetch timeseries via Tinkoff API due to: "+err.Error(), logger.ALERT)
		c.JSON(http.StatusInternalServerError, "Failed to fetch time series due to: " + err.Error())
	}
}
