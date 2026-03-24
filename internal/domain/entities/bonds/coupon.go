package bonds

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

type Coupon struct {
	Id              uuid.UUID
	Figi            string
	CouponDate      time.Time
	CouponNumber    int
	RecordDate      time.Time
	PerBondAmount   float64
	CouponType      CouponType
	CouponStartDate time.Time
	CouponEndDate   time.Time
	CouponPeriod    int
}

type CouponType int32

const (
	CouponType_COUPON_TYPE_UNSPECIFIED CouponType = 0 //Неопределенное значение.
	CouponType_COUPON_TYPE_CONSTANT    CouponType = 1 //Постоянный.
	CouponType_COUPON_TYPE_FLOATING    CouponType = 2 //Плавающий.
	CouponType_COUPON_TYPE_DISCOUNT    CouponType = 3 //Дисконт.
	CouponType_COUPON_TYPE_MORTGAGE    CouponType = 4 //Ипотечный.
	CouponType_COUPON_TYPE_FIX         CouponType = 5 //Фиксированный.
	CouponType_COUPON_TYPE_VARIABLE    CouponType = 6 //Переменный.
	CouponType_COUPON_TYPE_OTHER       CouponType = 7 //Прочее.
)

// Enum value maps for CouponType.
var (
	CouponType_name = map[int32]string{
		0: "COUPON_TYPE_UNSPECIFIED",
		1: "COUPON_TYPE_CONSTANT",
		2: "COUPON_TYPE_FLOATING",
		3: "COUPON_TYPE_DISCOUNT",
		4: "COUPON_TYPE_MORTGAGE",
		5: "COUPON_TYPE_FIX",
		6: "COUPON_TYPE_VARIABLE",
		7: "COUPON_TYPE_OTHER",
	}
	CouponType_value = map[string]int32{
		"COUPON_TYPE_UNSPECIFIED": 0,
		"COUPON_TYPE_CONSTANT":    1,
		"COUPON_TYPE_FLOATING":    2,
		"COUPON_TYPE_DISCOUNT":    3,
		"COUPON_TYPE_MORTGAGE":    4,
		"COUPON_TYPE_FIX":         5,
		"COUPON_TYPE_VARIABLE":    6,
		"COUPON_TYPE_OTHER":       7,
	}
)

func TotalCouponIncome(coupons []Coupon, includePastCoupons bool) float64 {
	if len(coupons) == 0 {
		return -1
	}

	totalCouponIncome := 0.0
	for _, coupon := range coupons {
		if includePastCoupons == false && coupon.CouponDate.Before(time.Now()) {
			continue
		}
		totalCouponIncome += coupon.PerBondAmount
	}
	return totalCouponIncome
}

func AccumulatedCouponIncome(bond Bond, toDate time.Time) (float64, error) {
	if len(bond.Coupons) == 0 {
		return -1, errors.New("Attempting to calculate the accumulated coupon income with no coupons for " + bond.Figi)
	}

	currentCoupon, err := findCurrentCouponForBond(bond)
	if err != nil {
		return -1, err
	}
	if currentCoupon.CouponPeriod <= 0 {
		return -1, errors.New("Coupon period for " + bond.Figi + " is invalid")
	}
	daysFractional := toDate.Sub(currentCoupon.CouponStartDate).Hours() / 24
	roundedDays := math.Trunc(daysFractional)+1
	daysElapsedSinceCouponStartDate := int(roundedDays)
	if daysElapsedSinceCouponStartDate == 0 {
		return 0, nil
	}

	couponAmountPerDay := currentCoupon.PerBondAmount / float64(currentCoupon.CouponPeriod)
	aci := couponAmountPerDay * float64(daysElapsedSinceCouponStartDate)
	roundedAci := math.Round(aci*100)/100
	aciInRub := roundedAci * 81.8763
	
	return aciInRub, nil
}

// TODO: Unit tests
func findCurrentCouponForBond(bond Bond) (Coupon, error) {
	if len(bond.Coupons) == 0 {
		return Coupon{}, errors.New("Attempting to find a coupon with missing coupons for " + bond.Figi)
	}

	for i, coupon := range bond.Coupons {
		if time.Now().After(coupon.CouponStartDate) && time.Now().Before(coupon.CouponEndDate) {
			return bond.Coupons[i], nil
		}
	}

	return Coupon{}, errors.New("Failed to find the current coupon for " + bond.Figi)
}
