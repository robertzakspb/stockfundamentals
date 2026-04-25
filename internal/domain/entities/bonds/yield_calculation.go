package bonds

import (
	"errors"
	"time"
)

func (b Bond) CalcYieldToMaturity(coupons []Coupon, marketPrice float64) (float64, error) {
	yield, err := calculateYield(b, coupons, marketPrice, time.Now(), b.MaturityDate, b.MaturityDate)
	if err != nil {
		return -1, err
	}

	return yield, nil
}

func (b Bond) CalcYieldToCallOption(coupons []Coupon, marketPrice float64) (float64, error) {
	if b.CallOptionExerciseDate.IsZero() {
		return -1, errors.New("Attempting to calculate a yield to call option for a bond without a call exercise date")
	}

	yield, err := calculateYield(b, coupons, marketPrice, time.Now(), b.CallOptionExerciseDate, b.CallOptionExerciseDate)
	if err != nil {
		return -1, err
	}

	return yield, nil
}

func calculateYield(b Bond, coupons []Coupon, marketPricePercentage float64, acquisitionDate, redemptionDate, latestCouponDate time.Time) (float64, error) {
	if len(coupons) == 0 {
		return -1, errors.New("Failed to calculate the yield due to missing coupons")
	}
	if !(coupons[0].CouponType == CouponType_COUPON_TYPE_FIX || coupons[0].CouponType == CouponType_COUPON_TYPE_CONSTANT) {
		return -1, errors.New("Unable to calculate the YTM for non-fixed and non-constant coupons")
	}
	if marketPricePercentage == 0 {
		yield := -1
		return float64(yield), errors.New("Unable to calculate the yield due to a missing market price")
	}

	holdingPeriod := redemptionDate.Sub(acquisitionDate).Hours() / 24

	marketPriceInCurrency := marketPricePercentage * b.NominalValue / 100

	marketPrice := marketPriceInCurrency + b.AccruedInterest

	tci := TotalCouponIncome(coupons, false, latestCouponDate)

	//Standard simple formula for the calculation of bond yields (coupon reinvestment is not assumed)
	yield := (b.NominalValue - marketPrice + tci) / marketPrice * 365 / holdingPeriod * 100
	return yield, nil
}

// Calculates the return realized on the bond given a market price, including coupons and redemption
// // Coupon reinvestment is not assumed
// func totalBondReturn(bond Bond, coupons []Coupon, marketPrice float64, acquisitionDate, redemptionDate time.Time) (float64, error) {
// 	if marketPrice == 0 {
// 		return -1, errors.New("Invalid market price")
// 	}

// 	futureCashflows := totalCouponIncome(coupons, false) + bond.NominalValue

// 	cumulativeReturn := futureCashflows/marketPrice - 1
// 	totalReturn := compoundinterest.CalcAnnualizedReturn(cumulativeReturn, acquisitionDate, redemptionDate)

// 	return totalReturn, nil
// }
