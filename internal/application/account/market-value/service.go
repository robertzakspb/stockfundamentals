package accountmvservice

import (
	"errors"
	"strconv"
	"time"

	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
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

	return accountMVs, nil
}

func AccountStockMarketValueGroupedByCurrency(accountId uuid.UUID, date time.Time) (map[string]accountmvdomain.AccountMarketValue, error) {
	accountPortfolio, err := portfolio.GetAccountPortfolio([]uuid.UUID{accountId})
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

	stockPortfolioMarketValue, currency, err := portfolio.CalculatePortfolioMarketValue(accountPortfolio, "RUB")
	if err != nil {
		return map[string]accountmvdomain.AccountMarketValue{}, err
	}

	mv := map[string]accountmvdomain.AccountMarketValue{"RUB": {
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
	figis := bondportfolio.GetLotFigis(bondLots)

	config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return accountmvdomain.AccountMarketValue{}, err
	}
	quotes, err := tquoteservice.FetchQuotesForFigis(figis, config)

	totalMarketValue := 0.0

	bonds := bondportfolio.GetLotBonds(bondLots)
	currencyPairs := bondservice.AllCurrencyPairsInBondList(bonds)
	fxRates, err := forexservice.GetExchangeRates(currencyPairs, date)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return accountmvdomain.AccountMarketValue{}, err
	}

	for _, quote := range quotes {
		foundQuote := false
		for _, lot := range bondLots {
			if lot.Figi == quote.Figi() {
				foundQuote = true

				fxRate := 1.0
				if lot.Bond.IsBondWithDifferentNominalCurrencyAndCurrency() {
					rate, found := forexservice.FindRate(lot.Bond.NominalCurrency, lot.Bond.Currency, fxRates)
					if !found {
						logger.Log("Failed to find an exchange rate for "+lot.Bond.NominalCurrency+"/"+lot.Bond.Currency+". Unable to calculate the market value for the bond "+lot.Bond.Isin, logger.ERROR)
						return accountmvdomain.AccountMarketValue{}, errors.New("Unable to calculate the market value for the bond due to the missing forex rate")
					}
					fxRate = rate.Rate
				}

				lotMarketValue := lot.MarketValue(quote.Quote(), fxRate)

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
