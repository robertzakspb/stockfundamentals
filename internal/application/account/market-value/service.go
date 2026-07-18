package accountmvservice

import (
	"errors"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"

	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
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
	stockMVs, err := GetAccountStockMarketValueGroupedByCurrency(accountId, date)
	if err != nil {
		return []accountmvdomain.AccountMarketValue{}, err
	}

	bondMVs, err := GetAccountBondMarketValueGroupedByCurrency(accountId, date)
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
