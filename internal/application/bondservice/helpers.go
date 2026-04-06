package bondservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func GetBondFigis(bondList *[]bonds.Bond) []string {
	figis := []string{}

	for _, bond := range *bondList {
		if bond.Figi == "" {
			continue
		}
		figis = append(figis, bond.Figi)
	}

	return figis
}

func GetOnlyBondsWithFixedOrConstantCoupons(bondList []bonds.Bond) []bonds.Bond {
	filteredBonds := []bonds.Bond{}
	for _, bond := range bondList {
		if len(bond.Coupons) == 0 {
			logger.Log("Attempting to find bonds with fixed or constant coupons for a bond with no coupons", logger.WARNING)
			continue
		}
		if bond.Coupons[0].CouponType == bonds.CouponType_COUPON_TYPE_CONSTANT || bond.Coupons[0].CouponType == bonds.CouponType_COUPON_TYPE_FIX {
			filteredBonds = append(filteredBonds, bond)
		}
	}
	return filteredBonds
}

func AllCurrencyPairsInBondList(bondList []bonds.Bond) []string {
	pairs := []string{}

	for _, bond := range bondList {
		if bond.Currency != bond.NominalCurrency {
			foundPair := false
			for _, pair := range pairs {
				if pair == bond.NominalCurrency+"/"+bond.Currency {
					foundPair = true
				}
			}
			if !foundPair {
				pairs = append(pairs, bond.NominalCurrency+"/"+bond.Currency)
			}
		}
	}
	return pairs
}

func MatchCouponsWithBonds(coupons []bonds.Coupon, bonds []bonds.Bond) []bonds.Bond {
	for _, coupon := range coupons {
		for i, b := range bonds {
			if coupon.Figi == b.Figi {
				bonds[i].Coupons = append(b.Coupons, coupon)
			}
		}
	}
	return bonds
}
