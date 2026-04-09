package accountmv

import (
	"errors"
	"strconv"

	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/dividend-calendar/market-value"
	accountmvdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/market-value"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
)

func GetAccountReturn(filters []ydbfilter.YdbFilter) (accountmvdomain.Return, error) {
	dbMarketValues, err := accountmvdb.GetAccountMarketValues(filters)
	if err != nil {
		return accountmvdomain.Return{}, err
	}
	if len(dbMarketValues) != 2 {
		return accountmvdomain.Return{}, errors.New("Expected two market values but got " + strconv.Itoa(len(dbMarketValues)))
	}

	marketValues := []accountmvdomain.AccountMarketValue{}
	for _, dbMarketValue := range dbMarketValues {
		marketValues = append(marketValues, mapAccountMarketValueDbModelToDomain(dbMarketValue))
	}

	totalReturn := accountmvdomain.CalculateAccountReturn(marketValues[0].AccountId, marketValues[0], marketValues[1])

	return totalReturn, nil
}
