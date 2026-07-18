package accountmvservice

import (
	"errors"
	"time"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/quoteservice"
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

func GetAccountBondMarketValueGroupedByCurrency(accountId uuid.UUID, date time.Time) (map[string]accountmvdomain.AccountMarketValue, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.UuidValue(accountId),
	}
	bondLots, err := bondportfolio.GetFilteredPositionLots([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return map[string]accountmvdomain.AccountMarketValue{}, err
	}
	if len(bondLots) == 0 {
		return map[string]accountmvdomain.AccountMarketValue{}, nil
	}

	bondLots, err = bondportfolio.PopulateLotsWithBonds(bondLots)
	if err != nil {
		return map[string]accountmvdomain.AccountMarketValue{}, err
	}

	lotsGroupedByNominalCurrency := bondportfolio.GroupByNominalCurrency(bondLots)
	marketValues := map[string]accountmvdomain.AccountMarketValue{}
	for currency, lots := range lotsGroupedByNominalCurrency {
		lotsMarketValue, err := CalculateBondLotsMarketValue(lots, date, currency)
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
			continue
		}
		marketValues[currency] = lotsMarketValue
	}
	return marketValues, nil
}

func CalculateBondLotsMarketValue(bondLots []bonds.BondLot, date time.Time, currency string) (accountmvdomain.AccountMarketValue, error) {
	if len(bondLots) == 0 {
		return accountmvdomain.AccountMarketValue{}, errors.New("Attempting to calculate the market value of 0 bonds")
	}

	figis := bondportfolio.GetLotFigis(bondLots)

	quotes, err := quoteservice.FetchStockQuotes(figis)

	totalMarketValue := 0.0

	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return accountmvdomain.AccountMarketValue{}, err
	}

	for _, quote := range quotes {
		foundQuote := false
		for _, lot := range bondLots {
			if lot.Figi == quote.Figi() {
				foundQuote = true
				lotMarketValue := lot.MarketValue(quote.Quote(), 1.0)
				totalMarketValue += lotMarketValue
			}
		}
		if !foundQuote {
			return accountmvdomain.AccountMarketValue{}, errors.New("Failed to find the quote for figi: " + quote.Figi() + ". Terminating the market value calculation")
		}
	}

	marketValue := accountmvdomain.AccountMarketValue{
		AccountId: bondLots[0].AccountId,
		Date:      date,
		Currency:  currency,
		EodValue:  totalMarketValue,
	}

	return marketValue, nil
}
