package portfolio

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/quotefetcher"
	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/portfoliodb"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

func UpdatePortfolio() error {
	portfolio, _ := GeStockPortfolio()
	return portfoliodb.UpdateLocalPortfolio(mapLotToDbLot(portfolio.Lots))
}

func GetAccountPortfolio(filters []ydbfilter.YdbFilter) (stockportfolio.Portfolio, error) {
	dbLots, err := portfoliodb.GetAccountPortfolio(filters)
	if err != nil {
		return stockportfolio.Portfolio{}, err
	}

	lots := []lot.Lot{}
	for _, lot := range dbLots {
		lots = append(lots, mapLotDbToLot(lot))
	}

	return stockportfolio.Portfolio{Lots: lots}, nil
}

// Returns the market value alongside the used currency
func CalculatePortfolioMarketValue(portfolio stockportfolio.Portfolio, currency string) (float64, string, error) {
	portfolio, err := PopulateLotsWithQuotes(portfolio)
	if err != nil {
		return -1, currency, err
	}

	currencies := portfolio.PositionCurrencies()
	currencyPairs := []string{}
	for _, positionCurrency := range currencies {
		if positionCurrency != currency {
			currencyPairs = append(currencyPairs, positionCurrency+"/"+currency)
		}
	}
	forexRates, err := forexservice.GetExchangeRates(currencyPairs, time.Now())
	if err != nil {
		return -1, currency, err
	}

	totalMarketValue := 0.0
	for _, position := range portfolio.Lots {
		mv, err := position.MarketValue()
		if err != nil {
			return -1, currency, err
		}
		//Case when the position's currency is already the target currency
		if position.Currency == currency {
			totalMarketValue += mv
			continue
		}
		//Case when the position's market value must be calculated in the target currency
		rate, found := forexservice.FindRate(position.Currency, currency, forexRates)
		if !found {
			return -1, currency, errors.New("Failed to find an exchange rate for " + position.Currency + "/" + currency)
		}
		totalMarketValue += mv * rate.Rate
	}

	return totalMarketValue, currency, nil
}

func PopulateLotsWithQuotes(portfolio portfolio.Portfolio) (portfolio.Portfolio, error) {
	positions := portfolio.UniquePositions()

	securities, err := security_master.GetSecuritiesFilteredByFigi(stockportfolio.LotFigis(positions))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return portfolio, err
	}

	positions, errorList := stockportfolio.MatchLotsWithStocks(positions, securities)
	if len(errorList) > 0 {
		logger.Log(errorList[0].Error(), logger.ERROR)
		return portfolio, err
	}

	entitySecurities := []entity.Security{}
	for _, s := range securities {
		if s.GetId() == "" {
			continue
		}
		entitySecurities = append(entitySecurities, entity.Security{
			Figi:   s.GetFigi(),
			ISIN:   s.GetIsin(),
			Ticker: s.GetTicker(),
			MIC:    s.GetMic(),
		})
	}

	quotes := quotefetcher.FetchQuotesFor(entitySecurities)

	//Fetching quotes through Tinkoff API, as MOEX does not return quotes for ETFs like LQDT
	if etfFigis := portfolio.GetEtfLotFigis(); len(etfFigis) > 0 {
		config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
		if err != nil {
			logger.Log("Failed to initialize the configuration file", logger.ALERT)
			return portfolio, errors.New("Failed to fetch quotes for ETFs in the portfolio fue to Tinkoff API configuration issues")
		}
		etfQuotes, err := tquoteservice.FetchQuotesForFigis(portfolio.GetEtfLotFigis(), config)
		if err != nil {
			return portfolio, errors.New("Failed to fetch quotes for ETFs in the portfolio")
		}

		for _, etfQuote := range etfQuotes {
			quotes = append(quotes, &etfQuote)
		}
	}

	if len(quotes) != len(positions) {
		return portfolio, errors.New("The portfolio has " + strconv.Itoa(len(positions)) + " unique positions, whereas only " + strconv.Itoa(len(quotes)) + " quotes has been fetched")
	}

	positions, err = stockportfolio.MatchLotsWithQuotes(positions, quotes)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}

	sort.Slice(positions, func(i, j int) bool {
		return positions[i].CurrentReturn() > positions[j].CurrentReturn()
	})

	return stockportfolio.Portfolio{Lots: positions}, nil
}

func PopulateLotSecurities(lots []lot.Lot) ([]lot.Lot, error) {
	securities, err := security_master.GetSecuritiesFilteredByFigi(stockportfolio.LotFigis(lots))
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return lots, err
	}

	lots, errorList := stockportfolio.MatchLotsWithStocks(lots, securities)
	if len(errorList) != 0 {
		for i := range errorList {
			logger.Log(errorList[i].Error(), logger.ERROR)
		}
		return lots, errorList[0]
	}

	return lots, nil
}
