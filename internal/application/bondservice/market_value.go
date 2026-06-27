package bondservice

import (
	"time"

	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func CalculateRubMarketValue(bondList []bonds.Bond, quotes []tquoteservice.TQuote) []bonds.Bond { 
	pairs := AllCurrencyPairsInBondList(bondList)

	rates, err := forexservice.GetExchangeRates(pairs, time.Now())
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}

	for _, quote := range quotes {
		for i := range bondList {
			if quote.Figi() == bondList[i].Figi {
				if bondList[i].NominalCurrency == "RUB" {
					bondList[i].MarketValueInRUB = bondList[i].MarketValue(quote.Quote(), 1.0)
					continue
				}
				fxRate, found := forexservice.FindRate(bondList[i].NominalCurrency, "RUB", rates)
				if !found {
					logger.Log("Failed to find the exchange rate for "+bondList[i].NominalCurrency+"/RUB", logger.ERROR)
				}
				bondList[i].MarketValueInRUB = bondList[i].MarketValue(quote.Quote(), fxRate.Rate)
			}
		}
	}

	return bondList
}
