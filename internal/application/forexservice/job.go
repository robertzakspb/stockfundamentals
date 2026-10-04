package forexservice

import (
	"errors"
	"fmt"
	"strings"
	"time"

	forexdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/forex"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
)

func StartImportForexRatesJob() {
	go ImportForexRates()
}

func ImportForexRates() {
	for _, currencyPair := range SupportedCurrencyPairs() {
		split := strings.Split(currencyPair, "/")
		cur1 := split[0]
		cur2 := split[1]
		FetchAndSaveCurrencyPairQuotes(string(cur1), string(cur2))
	}
}

func FetchAndSaveCurrencyPairQuotes(cur1, cur2 string) error {
	earliestDateInDb, latestDateInDb, err := forexdb.GetEarliestAndLatestDbRateFor(cur1, cur2)
	if err != nil {
		return err
	}

	needToSkipDatesAlreadyInDb := true
	if earliestDateInDb.IsZero() || latestDateInDb.IsZero() {
		needToSkipDatesAlreadyInDb = false
	}

	cbrRateLimit := time.Second / 2
	cbrThrottle := time.Tick(cbrRateLimit)

	nbsRateLimit := time.Second / 2
	nbsThrottle := time.Tick(nbsRateLimit)

	rates := []ForexRate{}
	targetDate := time.Now().Add(-time.Hour * 24 * 365)
	for {
		if targetDate.After(time.Now().AddDate(0, 0, 1)) {
			break
		}

		if needToSkipDatesAlreadyInDb {
			if !(targetDate.Before(earliestDateInDb) || targetDate.After(latestDateInDb)) {
				targetDate = targetDate.Add(time.Hour * 24)
				continue //If the rate for the date is in the DB, don't import it
			}
		}

		var rate ForexRate
		var err error

		switch cur2 {
		case "RUB":
			<-cbrThrottle
			rate, err = getCurrencyToRubRate(cur1, targetDate)
		case "RSD":
			<-nbsThrottle
			rate, err = fetchUsdToRsdRate(targetDate)
		default:
			return errors.New("Unsupported currency: " + cur2)
		}

		if err != nil {
			//The rate for the following day may or may not be provided; hence, a warning is sufficient. Otherwise, an error.
			if timehelpers.AreEqualDates(time.Now().AddDate(0, 0, 1), targetDate) {
				logger.Log(err.Error(), logger.WARNING)
			} else {
				logger.Log(err.Error(), logger.ERROR)
			}

			targetDate = targetDate.Add(time.Hour * 24)
			continue
		}
		rates = append(rates, rate)
		logger.Log("Fetched the rate for "+cur1+"/"+cur2+". Value: "+fmt.Sprint(rate)+" for "+targetDate.String(), logger.INFORMATION)

		targetDate = targetDate.Add(time.Hour * 24)
	}

	usdToEurRates, err := FetchUsdToEurRate(time.Now().Add(-time.Hour*24*365), time.Now())
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	} else {
		rates = append(rates, usdToEurRates...)
	}

	mappedDbModels := mapFxRatesToDbModel(rates)
	err = forexdb.SaveForexRates(mappedDbModels)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}
