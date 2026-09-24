package accountmvservice

import (
	"errors"
	"sync"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/quoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
)

// Assets Under Management
type AUM struct {
	TotalAum float64
	Currency string
	Date     time.Time
}

func GetTotalAssetsUnderManagement(currencies ...string) ([]AUM, error) {
	var stockLots []lot.Lot
	var bondLots []bonds.BondLot
	var stockQuotes []entity.SimpleQuote
	var bondQuotes []entity.BondQuote
	var rates []forexservice.ForexRate
	var wg sync.WaitGroup
	var err error

	//For the first parallel batch, we fetch stock and bond lots
	wg.Go(func() {
		stockLots, err = portfolio.GetAllStockLots()
		stockLots, err = portfolio.PopulateLotSecurities(stockLots)
	})
	wg.Go(func() {
		bondLots, err = bondportfolio.GetAllPositionLots()
		bondLots, err = bondportfolio.PopulateLotsWithBonds(bondLots)
	})
	wg.Wait()

	//For the second parallel batch, we fetch stock and bond quotes as well as the forex rates
	wg.Go(func() {
		stockSecurities := portfolio.ExtractSecuritiesFromLots(stockLots)
		bondFigis := bondportfolio.GetLotFigis(bondLots)
		stockQuotes, bondQuotes, err = quoteservice.FetchInternationalStockAndBondQuotes(stockSecurities, bondFigis)
	})
	wg.Go(func() {
		var currencyPairs []string
		for _, currency := range currencies {
			currencyPairs = append(currencyPairs, bonds.GetCurrencyPairs(currency, bondLots)...)
			currencyPairs = append(currencyPairs, lot.GetCurrencyPairs(currency, stockLots)...)
		}

		currencyPairs = stringhelpers.RemoveDuplicatesFrom(currencyPairs)
		rates, err = forexservice.GetExchangeRates(currencyPairs, time.Now())
	})
	wg.Wait()

	if err != nil {
		return []AUM{}, err
	}

	bondLots = bondportfolio.MatchLotsWithQuotes(bondLots, bondQuotes)
	stockLots, err = stockportfolio.MatchLotsWithQuotes(stockLots, stockQuotes)
	if err != nil {
		return []AUM{}, err
	}

	aums := []AUM{}
	for _, currency := range currencies {
		aum := 0.0
		for i := range bondLots {
			if bondLots[i].Bond.NominalCurrency == currency {
				aum += bondLots[i].MarketValue(bondLots[i].Bond.QuoteInPercentage, 1.0)
			} else {
				rate, found := forexservice.FindRate(bondLots[i].Bond.NominalCurrency, currency, rates)
				if !found {
					logger.Log("Failed to find the exchange rate for "+bondLots[i].Bond.NominalCurrency+"/"+currency, logger.ERROR)
					continue
				}
				aum += bondLots[i].MarketValue(bondLots[i].Bond.QuoteInPercentage, rate.Rate)
			}
		}

		for i := range stockLots {
			if stockLots[i].Currency == currency {
				aum += stockLots[i].Quote * stockLots[i].Quantity
			} else {
				rate, found := forexservice.FindRate(stockLots[i].Currency, currency, rates)
				if !found {
					return []AUM{}, errors.New("Failed to find the exchange rate for " + bondLots[i].Bond.Currency + "/" + currency)
				}
				aum += stockLots[i].Quote * rate.Rate * stockLots[i].Quantity
			}
		}
		aumStruct := AUM{
			TotalAum: aum,
			Currency: currency,
			Date:     time.Now(),
		}
		aums = append(aums, aumStruct)
	}

	return aums, nil
}
