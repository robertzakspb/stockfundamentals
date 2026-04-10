package accountmvservice

import (
	"errors"
	"strconv"
	"time"

	"github.com/compoundinvest/invest-core/quote/bondquote"
	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	accountmvdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/market-value"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"opensource.tbank.ru/invest/invest-go/investgo"
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

func CalculateAccountMarketValue(accountId uuid.UUID, date time.Time) (accountmvdomain.AccountMarketValue, error) {
	// stockMV := CalculateAccountStockMarketValue(accountId, date)
	bondMV, err := CalculateAccountBondMarketValue(accountId, date)
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}

	totalMV := bondMV.EodValue //+ stockMV.EodValue

	mv := accountmvdomain.AccountMarketValue{
		AccountId: accountId,
		Date:      date,
		Currency:  "RUB",
		EodValue:  totalMV,
	}
	return mv, nil
}

// func CalculateAccountStockMarketValue(accountId uuid.UUID, date time.Time) accountmvdomain.AccountMarketValue {

// }

func CalculateAccountBondMarketValue(accountId uuid.UUID, date time.Time) (accountmvdomain.AccountMarketValue, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.UuidValue(accountId),
	}
	bondLots, err := bondportfolio.GetFilteredPositionLots([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}

	bondLots, err = bondportfolio.PopulateLotsWithBonds(bondLots)
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}

	figis := bondportfolio.GetLotFigis(bondLots)

	config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return accountmvdomain.AccountMarketValue{}, err
	}

	quotes, err := bondquote.FetchQuotesForFigis(figis, config)

	totalMarketValue := 0.0

	for _, quote := range quotes {
		foundQuote := false
		for _, lot := range bondLots {
			if lot.Figi == quote.Figi() {
				foundQuote = true
				marketPriceInCurrency := quote.QuoteAsPercentage() * lot.Bond.NominalValue / 100
				lotMarketValue := lot.Quantity * marketPriceInCurrency
				totalMarketValue += lotMarketValue
			}
		}
		if !foundQuote {
			return accountmvdomain.AccountMarketValue{}, errors.New("Failed to find the quote for figi: " + quote.Figi() + ". Terminating the market value calculation")
		}
	}

	marketValue := accountmvdomain.AccountMarketValue{
		AccountId: accountId,
		Date:      date,
		Currency:  "RUB",
		EodValue:  totalMarketValue,
	}

	return marketValue, nil
}
