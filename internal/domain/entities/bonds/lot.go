package bonds

import (
	"time"

	"github.com/google/uuid"
)

type BondLot struct {
	Id                uuid.UUID
	Figi              string
	Isin              string
	OpeningDate       time.Time
	ModificationDate  time.Time
	AccountId         uuid.UUID
	Quantity          float64
	PricePerUnit      float64
	PricePerUnitInRUB float64
	AccruedInterest   float64
	Bond              Bond
}

func (lot BondLot) PricePerUnitPercentage(nominalValue float64) float64 {
	if nominalValue != 0 {
		return lot.PricePerUnit / nominalValue
	}
	return -1
}

func (lot BondLot) CouponPayoutForPosition(coupon Coupon) float64 {
	return coupon.PerBondAmount * lot.Quantity
}

func (lot BondLot) TotalPrincipalRedemption(bond Bond) float64 {
	return bond.NominalValue * lot.Quantity
}
