package bondservice

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"opensource.tbank.ru/invest/invest-go/investgo"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
)

// Optimized method that fetches all data asynchronously
func PopulateBondsWithCouponsAndCalculateYtm(bondList []bonds.Bond) []bonds.Bond {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []bonds.Bond{}
	}
	client, err := investgo.NewClient(context.TODO(), config, nil)
	if err != nil {
		return []bonds.Bond{}
	}

	figis := GetBondFigis(&bondList)

	wg := sync.WaitGroup{}

	var coupons []bonds.Coupon
	wg.Go(func() {
		coupons, err = GetCouponsByFigis(figis)
	})

	var quotes []tquoteservice.BondQuote
	var errorList []error
	wg.Go(func() {
		quotes, errorList = tquoteservice.GetBondPriceAndYield(client, figis)
		if len(errorList) > 0 {
			logger.LogErrors(errorList, logger.ERROR)
		}
	})

	currencyPairs := AllCurrencyPairsInBondList(bondList)
	rates := []forexservice.ForexRate{}
	wg.Go(func() {
		rates, err = forexservice.GetExchangeRates(currencyPairs, time.Now())
		if err != nil {
			logger.Log(err.Error(), logger.ERROR)
		}
	})

	wg.Wait()

	bondList = MatchCouponsWithBonds(coupons, bondList)
	bondList = CalculateYtmForBondsUsingQuotes(bondList, quotes)
	bondList = CalculateRubMarketValue(bondList, quotes, rates)

	sort.Slice(bondList, func(i, j int) bool {
		return bondList[i].SimpleYieldToMaturity > bondList[i].SimpleYieldToMaturity
	})

	return bondList
}

func CalculateYtmForBonds(bondList []bonds.Bond) []bonds.Bond {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []bonds.Bond{}
	}
	client, err := investgo.NewClient(context.TODO(), config, nil)
	if err != nil {
		return []bonds.Bond{}
	}

	quotes, errorList := tquoteservice.GetBondPriceAndYield(client, GetBondFigis(&bondList))
	if len(errorList) > 0 {
		logger.LogErrors(errorList, logger.ERROR)
	}

	bondsWithYtm := CalculateYtmForBondsUsingQuotes(bondList, quotes)
	return bondsWithYtm
}

func CalculateYtmForBondsUsingQuotes(bondList []bonds.Bond, quotes []tquoteservice.BondQuote) []bonds.Bond {
	for _, quote := range quotes {
		for i, b := range bondList {
			if quote.Ticker == b.Ticker {
				bondList[i].QuoteInPercentage = quote.QuoteAsPercentage
				bondList[i].YieldTomaturity = quote.YTM
				if b.HasCallOption() {
					yieldToCallOption, err := b.CalcSimpleYieldToCallOption(b.Coupons, quote.QuoteAsPercentage)
					if err != nil {
						logger.Log(err.Error(), logger.ERROR)
					}
					bondList[i].SimpleYieldToCallOption = yieldToCallOption
					continue
				}

				ytm, err := b.CalcSimpleYieldToMaturity(b.Coupons, quote.QuoteAsPercentage)
				if err != nil {
					logger.Log(err.Error(), logger.ERROR)
				}
				bondList[i].SimpleYieldToMaturity = ytm
			}
		}
	}

	sort.Slice(bondList, func(i, j int) bool {
		return bondList[i].YieldTomaturity > bondList[j].YieldTomaturity
	})

	return bondList
}
