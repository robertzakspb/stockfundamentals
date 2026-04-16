package portfolio

import (
	"errors"
	"strconv"
	"time"

	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/invest-core/quote/quotefetcher"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/portfoliodb"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

func UpdatePortfolio() error {
	portfolio, _ := GeMyPortfolio()
	return portfoliodb.UpdateLocalPortfolio(mapLotToDbLot(portfolio.Lots))
}

func GetAccountPortfolio(accountIDs uuid.UUIDs) (stockportfolio.Portfolio, error) {
	dbLots, err := portfoliodb.GetAccountPortfolio(accountIDs)
	if err != nil {
		return stockportfolio.Portfolio{}, err
	}
	lots := []lot.Lot{}
	for _, lot := range dbLots {
		lots = append(lots, mapLotDbToLot(lot))
	}

	return stockportfolio.Portfolio{Lots: lots}, nil
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

// Returns the market value alongside the used currency and a possible error
func CalculatePortfolioMarketValue(portfolio stockportfolio.Portfolio, currency string) (float64, string, error) {
	lotSecurities := []entity.Security{}
	for _, s := range stockportfolio.LotStocks(portfolio.Lots) {
		if s.GetId() == "" {
			continue
		}
		lotSecurities = append(lotSecurities, entity.Security{
			Figi:   s.GetFigi(),
			ISIN:   s.GetIsin(),
			Ticker: s.GetTicker(),
			MIC:    s.GetMic(),
		})
	}

	quotes := quotefetcher.FetchQuotesFor(lotSecurities)

	config, err := investgo.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return -1, currency, errors.New("Failed to fetch quotes for ETFs in the portfolio fue to Tinkoff API configuration issues")
	}
	etfQuotes, err := tquoteservice.FetchQuotesForFigis(portfolio.GetEtfLotFigis(), config)
	if err != nil {
		return -1, currency, errors.New("Failed to fetch quotes for ETFs in the portfolio")
	}
	//TODO: Need to refactor the core package to support fetching of quotes via Tinkoff API
	for _, etfQuote := range etfQuotes {
		quotes = append(quotes, &etfQuote)
	}

	uniquePositions := portfolio.UniquePositions()
	if len(quotes) != len(portfolio.UniquePositions()) {
		return -1, currency, errors.New("The portfolio has " + strconv.Itoa(len(uniquePositions)) + " positions, whereas only " + strconv.Itoa(len(quotes)) + " quotes has been fetched")
	}

	currencies := portfolio.PositionCurrencies()
	currencyPairs := []string{}
	for _, positionCurrency := range currencies {
		if positionCurrency != currency {
			currencyPairs = append(currencyPairs, positionCurrency+"/"+currency)
		}
	}
	forexRates, _ := forexservice.GetExchangeRates(currencyPairs, time.Now())

	totalMarketValue := 0.0
	for _, quote := range quotes {
		foundQuote := false
		for _, position := range uniquePositions {
			if quote.Figi() == position.Figi {
				foundQuote = true
				if position.Currency == currency {
					totalMarketValue += position.Quantity * quote.Quote()
				} else {
					foundRate := false
					for _, rate := range forexRates {
						if rate.Currency1 == forexservice.Currency(position.Currency) {
							foundRate = true
							totalMarketValue += position.Quantity * quote.Quote() * rate.Rate
						}
					}
					if !foundRate {
						return -1, currency, errors.New("Failed to find an exchange rate for " + position.Currency + "/" + currency)
					}
				}
			}
		}
		if !foundQuote {
			logger.Log("Failed to find a quote for figi "+quote.Figi(), logger.ERROR)
			return -1, currency, errors.New("Failed to find a quote for figi " + quote.Figi())
		}
	}

	return totalMarketValue, currency, nil
}
