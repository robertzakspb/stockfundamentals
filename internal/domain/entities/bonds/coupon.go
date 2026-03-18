package bonds

import (
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

func totalCouponIncome(coupons []Coupon) float64 {
	if len(coupons) == 0 {
		return -1
	}

	totalCouponIncome := 0.0
	for _, coupon := range coupons {
		totalCouponIncome += coupon.PerBondAmount
	}
	return totalCouponIncome
}
