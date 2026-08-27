package bondservice

import (
	"sort"
	"sync"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/application/market-data/quoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

// Optimized method that fetches all data asynchronously
func PopulateBondsWithCouponsAndCalculateYtm(bondList []bonds.Bond) []bonds.Bond {
	figis := ExtractBondFigis(&bondList)

	wg := sync.WaitGroup{}
	var err error

	var coupons []bonds.Coupon
	wg.Go(func() {
		coupons, err = GetCouponsByFigis(figis)
	})

	var quotes []entity.BondQuote
	wg.Go(func() {
		quotes, err = quoteservice.FetchBondQuotes(ExtractBondFigis(&bondList))
		if err != nil {
			logger.LogError(err, logger.ERROR)
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
	figis := make([]string, len(bondList))
	for i := range bondList {
		figis[i] = bondList[i].Figi
	}
	quotes, err := quoteservice.FetchBondQuotes(figis)
	if err != nil {
		logger.LogError(err, logger.ERROR)
		return bondList
	}

	bondsWithYtm := CalculateYtmForBondsUsingQuotes(bondList, quotes)
	return bondsWithYtm
}

func CalculateYtmForBondsUsingQuotes(bondList []bonds.Bond, quotes []entity.BondQuote) []bonds.Bond {
	for _, quote := range quotes {
		for i, b := range bondList {
			if quote.GetTicker() == b.Ticker {
				bondList[i].QuoteInPercentage = quote.GetQuoteAsPercentage()
				bondList[i].YieldTomaturity = quote.GetYtm()
				if b.HasCallOption() {
					yieldToCallOption, err := b.CalcSimpleYieldToCallOption(b.Coupons, quote.GetQuoteAsPercentage())
					if err != nil {
						logger.Log(err.Error(), logger.ERROR)
					}
					bondList[i].SimpleYieldToCallOption = yieldToCallOption
					continue
				}

				ytm, err := b.CalcSimpleYieldToMaturity(b.Coupons, quote.GetQuoteAsPercentage())
				if err != nil {
					logger.Log(err.Error(), logger.ERROR)
				}
				bondList[i].SimpleYieldToMaturity = ytm
				bondList[i].YieldTomaturity = quote.GetYtm()
			}
		}
	}

	sort.Slice(bondList, func(i, j int) bool {
		return bondList[i].YieldTomaturity > bondList[j].YieldTomaturity
	})

	return bondList
}
