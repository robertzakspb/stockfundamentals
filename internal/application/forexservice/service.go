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
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

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

	rateLimit := time.Second / 2
	throttle := time.Tick(rateLimit)

	rates := []ForexRate{}
	targetDate := time.Now().Add(-time.Hour * 24 * 365)
	for {
		if targetDate.After(time.Now()) {
			break
		}

		if needToSkipDatesAlreadyInDb {
			if !(targetDate.Before(earliestDateInDb) || targetDate.After(latestDateInDb)) {
				targetDate = targetDate.Add(time.Hour * 24)
				continue //If the rate for the date is in the DB, don't import it
			}
		}

		rate, err := getCurrencyToRubRate(cur1, targetDate)
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
			continue
		}
		rates = append(rates, ForexRate{
			Currency1: currencyName[cur1],
			Currency2: currencyName[cur2],
			Date:      targetDate,
			Rate:      rate,
		})
		logger.Log("Fetched the rate for "+cur1+"/"+cur2+". Value: "+fmt.Sprint(rate)+" for "+targetDate.String(), logger.INFORMATION)

		targetDate = targetDate.Add(time.Hour * 24)

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

func GetExchangeRates(currencyPairs []string, date time.Time) ([]ForexRate, error) {
	if len(currencyPairs) == 0 {
		return []ForexRate{}, nil
	}
	filters := []ydbfilter.YdbFilter{{
		YqlColumnName:  "date",
		Condition:      ydbfilter.Equal,
		ConditionValue: ydbhelper.ConvertToYdbDate(date),
	}}

	cur1s := []string{}
	cur2s := []string{}
	for _, pair := range currencyPairs {
		split := strings.Split(pair, "/")
		cur1, cur2 := strings.ToUpper(split[0]), strings.ToUpper(split[1])
		if cur2 != "RUB" {
			logger.Log("Skipping the currency pair "+cur1+"/"+cur2+" due to missing rates", logger.INFORMATION)
			continue
		}
		cur1s = append(cur1s, cur1)
		cur2s = append(cur2s, cur2)
	}

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

func FindRate(cur1, cur2 string, rates []ForexRate) (ForexRate, bool) {
	for _, rate := range rates {
		if string(rate.Currency1) == strings.ToUpper(cur1) && string(rate.Currency2) == strings.ToUpper(cur2) {
			return rate, true
		}
	}

	return ForexRate{}, false
}

// func GetExchangeRatesObsolete(currencyPairs []string, date time.Time) ([]ForexRate, error) {
// 	rates := []ForexRate{}

// 	for _, pair := range currencyPairs {
// 		split := strings.Split(pair, "/")
// 		cur1, cur2 := strings.ToUpper(split[0]), strings.ToUpper(split[1])
// 		if cur2 != "RUB" {
// 			logger.Log("Skipping the currency pair "+cur1+"/"+cur2+" due to missing rates", logger.INFORMATION)
// 			continue
// 		}
// 		rate, err := GetExchangeRate(cur1, cur2, date)
// 		if err != nil {
// 			logger.Log("Failed to get the forex rate for "+cur1+"/"+cur2, logger.ERROR)
// 			continue
// 		}
// 		rates = append(rates, rate)
// 	}

// 	return rates, nil
// }

// func GetExchangeRateForPair(currency1, currency2 string, dp ForexDataProvider) (float64, error) {
// 	if currency1 == currency2 {
// 		return 1, nil
// 	}

// 	usdToCur1, err := dp.GetExchangeRateUsdTo(currency1)
// 	usdToCur2, err := dp.GetExchangeRateUsdTo(currency2)
// 	if err != nil || usdToCur2 == 0 {
// 		return 0, fmt.Errorf("%s", "failed to find the exchange rate between "+string(USD)+"and "+currency2)
// 	}

// 	return usdToCur1 / usdToCur2, nil
// }

// func ConvertPriceToDifferentCurrency(price float64, priceCur, targetCur string, dp ForexDataProvider) (float64, error) {
// 	exRate, err := GetExchangeRateForPair(targetCur, priceCur, dp)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return price / exRate, nil
// }
