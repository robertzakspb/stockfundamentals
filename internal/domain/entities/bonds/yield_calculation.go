package bonds

import (
	"errors"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/compoundinterest"
)

// Calculates the bond's yield-to-maturity using the internal IRR formula
func (b *Bond) CalculateYTM(quote float64) error {
	cashflows, dates, err := b.FutureCashflowsWithDates(false, quote, time.Now())
	if err != nil {
		return err
	}

	irr, err := compoundinterest.InternalRateOfReturn(cashflows, dates)
	if err != nil {
		return err
	}

	b.YieldTomaturity = irr

	return nil
}

// Calculates the bond's yield-to-call-option using the internal IRR formula
func (b *Bond) CalculateYieldToCallOption(quote float64) error {
	cashflows, dates, err := b.FutureCashflowsWithDates(true, quote, time.Now())
	if err != nil {
		return err
	}

	irr, err := compoundinterest.InternalRateOfReturn(cashflows, dates)
	if err != nil {
		return err
	}

	b.YieldToCallOption = irr
	
	return nil
}

func (b Bond) CalcSimpleYieldToMaturity(coupons []Coupon, marketPrice float64) (float64, error) {
	yield, err := calculateSimpleYield(b, coupons, marketPrice, time.Now(), b.MaturityDate, b.MaturityDate)
	if err != nil {
		return -1, err
	}

	return yield, nil
}

func (b Bond) CalcSimpleYieldToCallOption(coupons []Coupon, marketPrice float64) (float64, error) {
	if b.CallOptionExerciseDate.IsZero() {
		return -1, errors.New("Attempting to calculate a yield to call option for a bond without a call exercise date")
	}

	yield, err := calculateSimpleYield(b, coupons, marketPrice, time.Now(), b.CallOptionExerciseDate, b.CallOptionExerciseDate)
	if err != nil {
		return -1, err
	}

	return yield, nil
}

func calculateSimpleYield(b Bond, coupons []Coupon, marketPricePercentage float64, acquisitionDate, redemptionDate, latestCouponDate time.Time) (float64, error) {
	if len(coupons) == 0 {
		return -1, errors.New("Failed to calculate the yield due to missing coupons")
	}
	// if !(coupons[0].CouponType == CouponType_COUPON_TYPE_FIX || coupons[0].CouponType == CouponType_COUPON_TYPE_CONSTANT) {
	// 	return -1, errors.New("Unable to calculate the YTM for non-fixed and non-constant coupons")
	// }
	if marketPricePercentage == 0 {
		return -1, errors.New("Unable to calculate the yield due to a missing market price")
	}

	holdingPeriodInDays := redemptionDate.Sub(acquisitionDate).Hours() / 24

	marketPriceInCurrency := marketPricePercentage * b.NominalValue / 100

	marketPrice := marketPriceInCurrency + b.AccruedInterest

	tci := TotalCouponIncome(coupons, false, latestCouponDate)

	//Standard simple formula for the calculation of bond yields (coupon reinvestment is not assumed)
	yield := (b.NominalValue - marketPrice + tci) / marketPrice * 365 / holdingPeriodInDays * 100
	return yield, nil
}
