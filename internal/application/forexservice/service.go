package forexservice

import (
	"errors"
	"strconv"
	"strings"
	"time"

	forexdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/forex"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

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

func GetLatestAvailableRates(currencyPairs []string) ([]ForexRate, error) {
	if len(currencyPairs) == 0 {
		return []ForexRate{}, nil
	}
	cur1s, cur2s := generateCurrency1AndCurrency2Slices(currencyPairs)
	filters := []ydbfilter.YdbFilter{
		{
			YqlColumnName:  "currency_1",
			Condition:      ydbfilter.Contains,
			ConditionValue: ydbhelper.ConvertStringsToYdbList(cur1s),
		},
		{
			YqlColumnName:  "currency_2",
			Condition:      ydbfilter.Contains,
			ConditionValue: ydbhelper.ConvertStringsToYdbList(cur2s),
		}}

	dbRates, err := forexdb.GetLatestAvailableRates(filters)
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
