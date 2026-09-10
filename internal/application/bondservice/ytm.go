package bondservice

import (
	"fmt"
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
	bondList = CalculateSimpleYtmForBondsUsingQuotes(bondList, quotes)
	bondList = CalculateRubMarketValue(bondList, quotes, rates)

	sort.Slice(bondList, func(i, j int) bool {
		return bondList[i].SimpleYieldToMaturity > bondList[i].SimpleYieldToMaturity
	})

	return bondList
}

func CalculateSimpleYtmForBonds(bondList []bonds.Bond) []bonds.Bond {
	figis := make([]string, len(bondList))
	for i := range bondList {
		figis[i] = bondList[i].Figi
	}
	quotes, err := quoteservice.FetchBondQuotes(figis)
	if err != nil {
		logger.LogError(err, logger.ERROR)
		return bondList
	}

	bondsWithYtm := CalculateSimpleYtmForBondsUsingQuotes(bondList, quotes)
	return bondsWithYtm
}

func CalculateSimpleYtmForBondsUsingQuotes(bondList []bonds.Bond, quotes []entity.BondQuote) []bonds.Bond {
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

// This function is used primarily for benchmarking -- it compares the internal bond YTMs to MOEX's YTM yields
func CompareYTMs() {
	bondList, err := GetRussianGovernmentBondsWithFixedOrConstantCoupon()
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}

	bondsCopy := make([]bonds.Bond, len(bondList))
	copy(bondsCopy, bondList)

	quotes, err := quoteservice.FetchBondQuotes(ExtractBondFigis(&bondList))
	if err != nil {
		logger.LogError(err, logger.ERROR)
		return
	}

	bondsCopy = CalculateBondYtmsUsingInternalIrrFormula(bondsCopy, quotes)

	for i := range bondsCopy {
		if bondsCopy[i].HasCallOption() {
			difference := bondList[i].YieldToCallOption - bondsCopy[i].YieldToCallOption*100
			fmt.Println(bondList[i].Name+". MOEX YTM: ", bondList[i].YieldToCallOption, ". Internal: ", bondsCopy[i].YieldToCallOption*100, ". Difference: ", difference, "%")
		} else {
			difference := bondList[i].YieldTomaturity - bondsCopy[i].YieldTomaturity*100
			fmt.Println(bondList[i].Name+". MOEX YTM: ", bondList[i].YieldTomaturity, ". Internal: ", bondsCopy[i].YieldTomaturity*100, ". Difference: ", difference, "%")
		}
	}
}

// Given a list of bonds and their corresponding percentage quotes, calculates the bonds' YTM or YTCO using the internal IRR formula
func CalculateBondYtmsUsingInternalIrrFormula(bondList []bonds.Bond, quotes []entity.BondQuote) []bonds.Bond {
	for i := range bondList {
		foundQuote := false
		for j := range quotes {
			if bondList[i].Ticker != quotes[j].GetTicker() {
				continue
			}
			foundQuote = true
			quoteInCurrency := bondList[i].MarketPriceInCurrency(quotes[j].GetQuoteAsPercentage())
			//Calculating either the current yield-to-maturity or the yield-to-call-option
			if bondList[i].HasCallOption() {
				err := bondList[i].CalculateYieldToCallOption(quoteInCurrency)
				if err != nil {
					break //If something went wrong, move on to the next bond
				}
			} else {
				err := bondList[i].CalculateYTM(quoteInCurrency)
				if err != nil {
					break //If something went wrong, move on to the next bond
				}
			}
			break //Once the yield are calculated for a bond, move on to the next bond
		}
		if !foundQuote {
			logger.Log("Failed to find a quote for bond "+bondList[i].Isin, logger.ERROR)
		}
	}
	return bondList
}
