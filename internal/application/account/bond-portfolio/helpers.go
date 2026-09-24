package bondportfolio

import (
	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func matchLotsWithBonds(lots []bonds.BondLot, bonds []bonds.Bond) []bonds.BondLot {
	if len(bonds) == 0 {
		return lots
	}

	for i, lot := range lots {
		for _, b := range bonds {
			if lot.Figi == b.Figi || lot.Isin == b.Isin {
				lots[i].Bond = b
			}
		}
	}

	return lots
}

func GetLotBonds(lots []bonds.BondLot) []bonds.Bond {
	bondList := []bonds.Bond{}

	for _, lot := range lots {
		if lot.Bond.Figi == "" && lot.Bond.Isin == "" && lot.Bond.NominalValue == 0 {
			continue
		}
		bondList = append(bondList, lot.Bond)
	}

	return bondList
}

func GetLotCoupons(lots []bonds.BondLot) []bonds.Coupon {
	couponList := []bonds.Coupon{}

	for _, lot := range lots {
		if lot.Bond.Figi == "" && lot.Bond.Isin == "" && lot.Bond.NominalValue == 0 {
			continue
		}
		for _, coupon := range lot.Bond.Coupons {
			couponList = append(couponList, coupon)
		}

	}

	return couponList
}

func GroupByNominalCurrency(lots []bonds.BondLot) map[string][]bonds.BondLot {
	groupedLots := map[string][]bonds.BondLot{}

	for i := range lots {
		groupedLots[lots[i].Bond.NominalCurrency] = append(groupedLots[lots[i].Bond.NominalCurrency], lots[i])
	}

	return groupedLots
}

func GetLotFigis(lots []bonds.BondLot) []string {
	figis := []string{}
	for _, lot := range lots {
		if lot.Figi != "" {
			figis = append(figis, lot.Figi)
		}
	}
	return figis
}

func GetLotIsins(lots []bonds.BondLot) []string {
	isins := []string{}
	for _, lot := range lots {
		if lot.Isin != "" {
			isins = append(isins, lot.Isin)
		}
	}
	return isins
}

func MatchLotsWithQuotes(lots []bonds.BondLot, quotes []entity.BondQuote) []bonds.BondLot {
	for i := range lots {
		foundQuote := false
		for j := range quotes {
			if quotes[j].GetTicker() != lots[i].Bond.Ticker {
				continue
			}
			foundQuote = true
			lots[i].Bond.QuoteInPercentage = quotes[j].GetQuoteAsPercentage()
		}
		if !foundQuote {
			logger.Log("Failed to find a quote for bond lot "+lots[i].Isin, logger.ERROR)
		}
	}
	return lots
}
