package forexservice

import (
	"errors"
	"strings"
	"time"

	forexdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/forex"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func StartFxRateImportJob() {
	var requiredCurrencyPairs = []string{"USD/RUB", "EUR/RUB"}
	for _, pair := range requiredCurrencyPairs {
		split := strings.Split(pair, "/")
		cur1, cur2 := split[0], split[1]
		FetchAndSaveCurrencyPairQuotes(cur1, cur2)
	}
}

func FetchAndSaveCurrencyPairQuotes(cur1, cur2 string) error {
	if cur2 != "RUB" {
		return errors.New("Unable to fetch pair quotes for non-Ruble currencies due to a missing implementation")
	}

	earliestDateInDb, latestDateInDb, err := forexdb.GetEarliestAndLatestDbRateFor(cur1, cur2)
	if err != nil {
		return err
	}

	needToSkipDatesAlreadyInDb := true
	if earliestDateInDb.IsZero() || latestDateInDb.IsZero() {
		needToSkipDatesAlreadyInDb = false
	}

	rateLimit := time.Second / 2 //To comply with the Tinkoff API rate limits
	throttle := time.Tick(rateLimit)

	rates := []ForexRate{}
	targetDate := time.Now().Add(-time.Hour * 24 * 365)
	for {
		if targetDate.After(time.Now()) {
			break
		}

		if needToSkipDatesAlreadyInDb {
			if !(targetDate.Before(earliestDateInDb) || targetDate.After(latestDateInDb)) {
				continue //If the rate for the date is in the DB, don't import it
			}
		}

		rate, err := GetCurrencyToRubRate("USD", targetDate)
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
			continue
		}
		rates = append(rates, ForexRate{
			Currency1: cur1,
			Currency2: cur2,
			Date:      targetDate,
			Rate:      rate,
		})

		targetDate.Add(time.Hour * 24)

		<-throttle
	}

	mappedDbModels := mapFxRatesToDbModel(rates)
	err = forexdb.SaveForexRates(mappedDbModels)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}
