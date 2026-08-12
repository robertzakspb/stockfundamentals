package forexservice

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	forexdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/forex"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

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

	mappedDbModels := mapFxRatesToDbModel(rates)
	err = forexdb.SaveForexRates(mappedDbModels)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	return nil
}

func GetExchangeRates(currencyPairs []string, date time.Time) ([]ForexRate, error) {
	if len(currencyPairs) == 0 {
		return []ForexRate{}, nil
	}

	cur1s, cur2s := generateCurrency1AndCurrency2Slices(currencyPairs)

	filters := []ydbfilter.YdbFilter{{
		YqlColumnName:  "date",
		Condition:      ydbfilter.Equal,
		ConditionValue: ydbhelper.ConvertToYdbDate(date),
	}}
	filters = append(filters, ydbfilter.YdbFilter{
		YqlColumnName:  "currency_1",
		Condition:      ydbfilter.Contains,
		ConditionValue: ydbhelper.ConvertStringsToYdbList(cur1s),
	})
	filters = append(filters, ydbfilter.YdbFilter{
		YqlColumnName:  "currency_2",
		Condition:      ydbfilter.Contains,
		ConditionValue: ydbhelper.ConvertStringsToYdbList(cur2s),
	})

	dbRates, err := forexdb.GetAllFxRates(filters)
	if err != nil {
		return []ForexRate{}, err
	}

	rates := mapDbModelsToDomain(dbRates)

	rates, err = collapseRatesIntoTargetCrossRates(currencyPairs, rates)
	if err != nil {
		return rates, err
	}

	return rates, nil
}

func GetExchangeRate(cur1, cur2 string, date time.Time) (ForexRate, error) {
	filters := []ydbfilter.YdbFilter{
		{
			YqlColumnName:  "date",
			Condition:      ydbfilter.Equal,
			ConditionValue: ydbhelper.ConvertToYdbDate(date),
		}, {
			YqlColumnName:  "currency_1",
			Condition:      ydbfilter.Equal,
			ConditionValue: types.TextValue(strings.ToUpper(cur1)),
		}, {
			YqlColumnName:  "currency_2",
			Condition:      ydbfilter.Equal,
			ConditionValue: types.TextValue(strings.ToUpper(cur2)),
		},
	}
	rates, err := forexdb.GetAllFxRates(filters)
	if err != nil {
		return ForexRate{}, err
	}
	if len(rates) == 0 || len(rates) > 1 {
		return ForexRate{}, errors.New("Invalid number of forex rates retrieved from the database: " + strconv.Itoa(len(rates)))
	}

	return mapDbModelsToDomain(rates)[0], nil
}

func GetFilteredExchangeRates(filters []ydbfilter.YdbFilter) ([]ForexRate, error) {
	rates, err := forexdb.GetAllFxRates(filters)
	if err != nil {
		return []ForexRate{}, err
	}
	if len(rates) == 0 {
		return []ForexRate{}, errors.New("Retrieved 0 forex rates from the database")
	}

	mappedRates := mapDbModelsToDomain(rates)

	return mappedRates, nil
}
