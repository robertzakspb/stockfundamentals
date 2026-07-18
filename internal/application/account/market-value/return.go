package accountmvservice

import (
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	accountmvdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/market-value"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
)

func GetAccountReturn(filters []ydbfilter.YdbFilter, currency string) (accountmvdomain.Return, error) {
	dbMarketValues, err := accountmvdb.GetAccountMarketValues(filters)
	if err != nil {
		return accountmvdomain.Return{}, err
	}

	marketValues := []accountmvdomain.AccountMarketValue{}
	for _, dbMarketValue := range dbMarketValues {
		marketValues = append(marketValues, mapAccountMarketValueDbModelToDomain(dbMarketValue))
	}

	startingDate := marketValues[0].Date
	startingDateMVs := []accountmvdomain.AccountMarketValue{}
	endingDateMVs := []accountmvdomain.AccountMarketValue{}

	for _, mv := range marketValues {
		if timehelpers.AreEqualDates(mv.Date, startingDate) {
			startingDateMVs = append(startingDateMVs, mv)
		} else {
			endingDateMVs = append(endingDateMVs, mv)
		}
	}

	startingDateMV, err := ConvertAccountMVsToCurrency(startingDateMVs, currency)
	if err != nil {
		return accountmvdomain.Return{}, err
	}
	endingDateMV, err := ConvertAccountMVsToCurrency(endingDateMVs, currency)
	if err != nil {
		return accountmvdomain.Return{}, err
	}

	totalReturn := accountmvdomain.CalculateAccountReturn(marketValues[0].AccountId, startingDateMV, endingDateMV)

	return totalReturn, nil
}
