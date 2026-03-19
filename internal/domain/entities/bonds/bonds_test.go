package bonds

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test__calculateYield(t *testing.T) {
	acquisitionDate, err := time.Parse("2006-01-02", "2026-03-19")
	test.AssertEqual(t, err, nil)
	maturityDate, err := time.Parse("2006-01-02", "2027-10-06")
	test.AssertEqual(t, err, nil)

	bond := Bond{
		MaturityDate:    maturityDate,
		NominalValue:    1000,
		NominalCurrency: "RUB",
	}

	coupons := []Coupon{
		{
			CouponType:    CouponType_COUPON_TYPE_CONSTANT,
			PerBondAmount: 29.92,
			CouponDate:    time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			CouponType:    CouponType_COUPON_TYPE_CONSTANT,
			PerBondAmount: 29.92,
			CouponDate:    time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			CouponType:    CouponType_COUPON_TYPE_CONSTANT,
			PerBondAmount: 29.92,
			CouponDate:    time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			CouponType:    CouponType_COUPON_TYPE_CONSTANT,
			PerBondAmount: 29.92,
			CouponDate:    time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
		},
	}

	yield, err := calculateYield(bond, coupons, 906.63+26.8, acquisitionDate, maturityDate)
	test.AssertEqual(t, err, nil)

	test.AssertEqualFloat(t, 12.87, yield, 0.01)
}
