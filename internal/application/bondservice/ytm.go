package bondservice

import (
	"sort"
	"sync"

	"github.com/compoundinvest/invest-core/quote/tquoteservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
)

// Optimized method that fetches all data asynchronously
func PopulateBondsWithCouponsAndCalculateYtm(bondList []bonds.Bond) []bonds.Bond {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []bonds.Bond{}
	}

	figis := GetBondFigis(&bondList)

	wg := sync.WaitGroup{}

	var coupons []bonds.Coupon
	wg.Go(func() {
		coupons, err = GetCouponsByFigis(figis)
	})

	var quotes []tquoteservice.TQuote
	wg.Go(func() {
		quotes, err = tquoteservice.FetchQuotesForFigis(figis, config)
	})

	wg.Wait()

	bondList = MatchCouponsWithBonds(coupons, bondList)
	bondList = CalculateYtmForBondsUsingQuotes(bondList, quotes)

	sort.Slice(bondList, func(i, j int) bool {
		return bondList[i].YieldToMaturity > bondList[i].YieldToMaturity
	})

	return bondList
}

func CalculateYtmForBonds(bondList []bonds.Bond) []bonds.Bond {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []bonds.Bond{}
	}

	quotes, err := tquoteservice.FetchQuotesForFigis(GetBondFigis(&bondList), config)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}

	bondsWithYtm := CalculateYtmForBondsUsingQuotes(bondList, quotes)
	return bondsWithYtm
}

func CalculateYtmForBondsUsingQuotes(bondList []bonds.Bond, quotes []tquoteservice.TQuote) []bonds.Bond {
	for _, quote := range quotes {
		for i, b := range bondList {
			if quote.Figi() == b.Figi {
				ytm, err := b.CalcYieldToMaturity(b.Coupons, quote.Quote())
				if err != nil {
					logger.Log(err.Error(), logger.ERROR)
					continue
				}
				bondList[i].YieldToMaturity = ytm

				if b.HasCallOption() {
					yieldToCallOption, err := b.CalcYieldToCallOption(b.Coupons, quote.Quote())
					if err != nil {
						logger.Log(err.Error(), logger.ERROR)
						continue
					}
					bondList[i].YieldToCallOption = yieldToCallOption
				}
			}
		}
	}
	return bondList
}
