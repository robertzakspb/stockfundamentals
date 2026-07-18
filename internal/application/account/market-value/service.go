package accountmvservice

import (
	"errors"
	"time"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/quoteservice"

	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

func ConvertAccountMVsToCurrency(MVs []accountmvdomain.AccountMarketValue, currency string) (accountmvdomain.AccountMarketValue, error) {
	if len(MVs) == 0 {
		return accountmvdomain.AccountMarketValue{}, errors.New("Zero market values were provided")
	}

	totalMV := 0.0
	currencyPairs := marketValueCurrencyPairs(currency, MVs)

	rates, err := forexservice.GetExchangeRates(currencyPairs, MVs[0].Date)
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}

	for _, mv := range MVs {
		if mv.Currency == currency {
			totalMV += mv.EodValue
		} else {
			rate, found := forexservice.FindRate(mv.Currency, currency, rates)
			if !found {
				rate, found = forexservice.FindRate(currency, mv.Currency, rates)
				if !found {
					return accountmvdomain.AccountMarketValue{}, errors.New("Failed to find the rate for " + mv.Currency + "/" + currency)
				}
				if rate.Rate != 0 {
					rate.Rate = 1 / rate.Rate
				}
			}
			adjustedMV := mv.EodValue * rate.Rate
			totalMV += adjustedMV
		}
	}

	return accountmvdomain.AccountMarketValue{
		AccountId: MVs[0].AccountId,
		Date:      MVs[0].Date,
		EodValue:  totalMV,
		Currency:  currency,
	}, nil
}

func CalculateAccountMarketValue(accountId uuid.UUID, date time.Time) ([]accountmvdomain.AccountMarketValue, error) {
	stockMVs, err := AccountStockMarketValueGroupedByCurrency(accountId, date)
	if err != nil {
		return []accountmvdomain.AccountMarketValue{}, err
	}

	bondMVs, err := AccountBondMarketValueGroupedByCurrency(accountId, date)
	if err != nil {
		return []accountmvdomain.AccountMarketValue{}, err
	}

	currencies := ExtractMarketValueCurrencies(stockMVs, bondMVs)

	accountMVs := []accountmvdomain.AccountMarketValue{}
	for i := range currencies {
		stockMV, foundStockMV := stockMVs[currencies[i]]
		bondMV, foundBondMV := bondMVs[currencies[i]]
		if foundStockMV && foundBondMV {
			accountMVs = append(accountMVs, accountmvdomain.AccountMarketValue{
				AccountId: stockMV.AccountId,
				Date:      stockMV.Date,
				Currency:  stockMV.Currency,
				EodValue:  stockMV.EodValue + bondMV.EodValue,
			})
		}
		if foundStockMV && !foundBondMV {
			accountMVs = append(accountMVs, accountmvdomain.AccountMarketValue{
				AccountId: stockMV.AccountId,
				Date:      stockMV.Date,
				Currency:  stockMV.Currency,
				EodValue:  stockMV.EodValue,
			})
		}
		if !foundStockMV && foundBondMV {
			accountMVs = append(accountMVs, accountmvdomain.AccountMarketValue{
				AccountId: bondMV.AccountId,
				Date:      bondMV.Date,
				Currency:  bondMV.Currency,
				EodValue:  bondMV.EodValue,
			})
		}
	}

	if len(accountMVs) == 0 {
		logger.Log("The account market value for account "+accountId.String()+" is 0.", logger.ALERT)
	}

	return accountMVs, nil
}

func AccountStockMarketValueGroupedByCurrency(accountId uuid.UUID, date time.Time) (map[string]accountmvdomain.AccountMarketValue, error) {
	accountFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.UuidValue(accountId),
	}
	accountPortfolio, err := portfolio.GetAccountPortfolio([]ydbfilter.YdbFilter{accountFilter})
	if err != nil {
		return map[string]accountmvdomain.AccountMarketValue{}, err
	}
	if len(accountPortfolio.Lots) == 0 {
		return map[string]accountmvdomain.AccountMarketValue{}, nil
	}

	accountPortfolio.Lots, err = portfolio.PopulateLotSecurities(accountPortfolio.Lots)
	if err != nil {
		return map[string]accountmvdomain.AccountMarketValue{}, err
	}
	if len(accountPortfolio.Lots) == 0 {
		return map[string]accountmvdomain.AccountMarketValue{}, nil
	}

	stockPortfolioMarketValue, currency, err := portfolio.CalculatePortfolioMarketValue(accountPortfolio, accountPortfolio.Lots[0].Currency)
	if err != nil {
		return map[string]accountmvdomain.AccountMarketValue{}, err
	}

	mv := map[string]accountmvdomain.AccountMarketValue{currency: {
		AccountId: accountId,
		Date:      date,
		Currency:  currency,
		EodValue:  stockPortfolioMarketValue,
	}}
	return mv, nil
}

func AccountBondMarketValueGroupedByCurrency(accountId uuid.UUID, date time.Time) (map[string]accountmvdomain.AccountMarketValue, error) {
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
