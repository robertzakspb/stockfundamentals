package accountmvservice

import (
	"errors"
	"strconv"
	"time"

	"github.com/compoundinvest/invest-core/quote/bondquote"
	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
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

func CalculateAccountMarketValue(accountId uuid.UUID, date time.Time, currency string) (accountmvdomain.AccountMarketValue, error) {
	stockMV, err := CalculateAccountStockMarketValue(accountId, date, currency)
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}
	bondMV, err := CalculateAccountBondMarketValue(accountId, date, currency)
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}

	totalMV := bondMV.EodValue + stockMV.EodValue

	mv := accountmvdomain.AccountMarketValue{
		AccountId: accountId,
		Date:      date,
		Currency:  "RUB",
		EodValue:  totalMV,
	}
	return mv, nil
}

func CalculateAccountStockMarketValue(accountId uuid.UUID, date time.Time, currency string) (accountmvdomain.AccountMarketValue, error) {
	accountPortfolio, err := portfolio.GetAccountPortfolio([]uuid.UUID{accountId})
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}
	if len(accountPortfolio.Lots) == 0 {
		mv := accountmvdomain.AccountMarketValue{
			AccountId: accountId,
			Date:      date,
			Currency:  currency,
			EodValue:  0,
		}
		return mv, nil
	}

	accountPortfolio.Lots, err = portfolio.PopulateLotSecurities(accountPortfolio.Lots)
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}

	stockPortfolioMarketValue, currency, err := portfolio.CalculatePortfolioMarketValue(accountPortfolio, currency)
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}

	mv := accountmvdomain.AccountMarketValue{
		AccountId: accountId,
		Date:      date,
		Currency:  currency,
		EodValue:  stockPortfolioMarketValue,
	}
	return mv, nil
}

func CalculateAccountBondMarketValue(accountId uuid.UUID, date time.Time, currency string) (accountmvdomain.AccountMarketValue, error) {
	filter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.UuidValue(accountId),
	}
	bondLots, err := bondportfolio.GetFilteredPositionLots([]ydbfilter.YdbFilter{filter})
	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}
	if len(bondLots) == 0 {
		marketValue := accountmvdomain.AccountMarketValue{
			AccountId: accountId,
			Date:      date,
			Currency:  currency,
			EodValue:  0,
		}
		return marketValue, nil
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

	bonds := bondportfolio.GetLotBonds(bondLots)
	currencyPairs := bondservice.AllCurrencyPairsInBondList(bonds)
	fxRates, err := forexservice.GetExchangeRates(currencyPairs, time.Now())
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
					}
					fxRate = rate.Rate
				}

				lotMarketValue := lot.MarketValue(quote.QuoteAsPercentage(), fxRate)

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
		Currency:  currency,
		EodValue:  totalMarketValue,
	}

	return marketValue, nil
}
