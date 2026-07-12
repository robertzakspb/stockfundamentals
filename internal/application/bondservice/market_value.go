package bondservice

import (
	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func CalculateRubMarketValue(bondList []bonds.Bond, quotes []entity.BondQuote, rates []forexservice.ForexRate) []bonds.Bond {
	for _, quote := range quotes {
		for i := range bondList {
			if quote.GetTicker() == bondList[i].Ticker {
				if bondList[i].NominalCurrency == "RUB" {
					bondList[i].MarketValueInRUB = bondList[i].MarketValue(quote.GetQuoteAsPercentage(), 1.0)
					continue
				}
				fxRate, found := forexservice.FindRate(bondList[i].NominalCurrency, "RUB", rates)
				if !found {
					logger.Log("Failed to find the exchange rate for "+bondList[i].NominalCurrency+"/RUB", logger.ERROR)
				}
				bondList[i].MarketValueInRUB = bondList[i].MarketValue(quote.GetQuoteAsPercentage(), fxRate.Rate)
			}
		}
	}

	return bondList
}
