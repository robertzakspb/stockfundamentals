package bonds

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_findCurrentCouponForBond_NoCoupons(t *testing.T) {
	bond := Bond{}

	_, err := findCurrentCouponForBond(bond)

	test.AssertError(t, err)
}

func Test_findCurrentCouponForBond_OnlyPastCoupon(t *testing.T) {
	bond := Bond{}

	bond.Coupons = append(bond.Coupons, Coupon{
		CouponStartDate: time.Now().Add(-time.Hour * 24),
		CouponEndDate:   time.Now().Add(-time.Hour * 12),
	})

	_, err := findCurrentCouponForBond(bond)

	test.AssertError(t, err)
}

func Test_findCurrentCouponForBond(t *testing.T) {
	bond := Bond{}
	bond.Coupons = append(bond.Coupons, Coupon{
		Figi:            "test_figi",
		CouponStartDate: time.Now(),
		CouponEndDate:   time.Now().Add(time.Hour * 24 * 30),
	})
	bond.Coupons = append(bond.Coupons, Coupon{
		CouponStartDate: time.Now().Add(-time.Hour * 24),
		CouponEndDate:   time.Now().Add(-time.Hour * 12),
	})

	coupon, err := findCurrentCouponForBond(bond)

	test.AssertNoError(t, err)
	test.AssertEqualStrings(t, coupon.Figi, "test_figi")
}

func Test_TotalCouponIncome_NoCoupons(t *testing.T) {
	coupons := []Coupon{}

	tci := TotalCouponIncome(coupons, false, time.Now())

	test.AssertEqual(t, -1, tci)
}

func Test_TotalCouponIncome_PastCouponsAreDisregarded(t *testing.T) {
	coupons := []Coupon{
		{
			PerBondAmount:   45,
			CouponStartDate: time.Now(),
			CouponEndDate:   time.Now().Add(time.Hour * 24 * 30),
			CouponDate:      time.Now().Add(time.Hour * 24 * 30),
		},
		{
			PerBondAmount:   70,
			CouponStartDate: time.Now().Add(-time.Hour * 24),
			CouponEndDate:   time.Now().Add(-time.Hour * 12),
			CouponDate:      time.Now().Add(-time.Hour * 12),
		},
	}

	tci := TotalCouponIncome(coupons, false, time.Time{})

	test.AssertEqual(t, 45, tci)
}

func Test_TotalCouponIncome_DisregardCouponsPastCertainDate(t *testing.T) {
	coupons := []Coupon{
		{
			PerBondAmount:   45,
			CouponStartDate: time.Now(),
			CouponEndDate:   time.Now().Add(time.Hour * 24 * 30),
			CouponDate:      time.Now().Add(time.Hour * 24 * 30),
		},
		{
			PerBondAmount:   70,
			CouponStartDate: time.Now().Add(-time.Hour * 24),
			CouponEndDate:   time.Now().Add(-time.Hour * 12),
			CouponDate:      time.Now().Add(-time.Hour * 12),
		},
		{
			PerBondAmount:   25,
			CouponStartDate: time.Now().Add(-time.Hour * 24 * 30),
			CouponEndDate:   time.Now().Add(-time.Hour * 24 * 15),
			CouponDate:      time.Now().Add(-time.Hour * 24 * 15),
		},
	}

	tci := TotalCouponIncome(coupons, true, time.Now())

	test.AssertEqual(t, 95, tci)
}

func Test_TotalCouponIncome(t *testing.T) {
	coupons := []Coupon{
		{
			PerBondAmount:   45,
			CouponStartDate: time.Now(),
			CouponEndDate:   time.Now().Add(time.Hour * 24 * 30),
			CouponDate:      time.Now().Add(time.Hour * 24 * 30),
		},
		{
			PerBondAmount:   70,
			CouponStartDate: time.Now().Add(-time.Hour * 24),
			CouponEndDate:   time.Now().Add(-time.Hour * 12),
			CouponDate:      time.Now().Add(-time.Hour * 12),
		},
		{
			PerBondAmount:   25,
			CouponStartDate: time.Now().Add(-time.Hour * 24 * 30),
			CouponEndDate:   time.Now().Add(-time.Hour * 24 * 15),
			CouponDate:      time.Now().Add(-time.Hour * 24 * 15),
		},
	}

	tci := TotalCouponIncome(coupons, true, time.Time{})

	test.AssertEqual(t, 140, tci)
}

func Test_AccruedInterest_NoCoupons(t *testing.T) {
	bond := Bond{}

	_, err := AccruedInterest(bond, time.Now(), 40)

	test.AssertError(t, err)
}

func Test_AccruedInterest_MissingCurrentCoupon(t *testing.T) {
	bond := Bond{}
	bond.Coupons = append(bond.Coupons, Coupon{
		PerBondAmount:   25,
		CouponStartDate: time.Now().Add(-time.Hour * 24 * 30),
		CouponEndDate:   time.Now().Add(-time.Hour * 24 * 15),
		CouponDate:      time.Now().Add(-time.Hour * 24 * 15),
	})

	_, err := AccruedInterest(bond, time.Now(), 40)

	test.AssertError(t, err)
}

func Test_AccruedInterest_ZeroAccruedInterest(t *testing.T) {
	bond := Bond{}
	bond.Coupons = append(bond.Coupons, Coupon{
		PerBondAmount:   25,
		CouponStartDate: time.Now(),
		CouponPeriod:    10,
		CouponEndDate:   time.Now().Add(time.Hour * 24 * 15),
		CouponDate:      time.Now().Add(time.Hour * 24 * 15),
	})

	_, err := AccruedInterest(bond, time.Now(), 40)

	test.AssertNoError(t, err)
}

func Test_AccruedInterest(t *testing.T) {
	bond := Bond{}
	bond.Coupons = []Coupon{
		{
			PerBondAmount:   45,
			CouponStartDate: time.Now(),
			CouponPeriod:    10,
			CouponEndDate:   time.Now().Add(time.Hour * 24 * 30),
			CouponDate:      time.Now().Add(time.Hour * 24 * 30),
		},
		{
			PerBondAmount:   70,
			CouponPeriod:    10,
			CouponStartDate: time.Now().Add(-time.Hour * 24),
			CouponEndDate:   time.Now().Add(-time.Hour * 12),
			CouponDate:      time.Now().Add(-time.Hour * 12),
		},
		{
			PerBondAmount:   25,
			CouponPeriod:    10,
			CouponStartDate: time.Now().Add(-time.Hour * 24 * 30),
			CouponEndDate:   time.Now().Add(-time.Hour * 24 * 15),
			CouponDate:      time.Now().Add(-time.Hour * 24 * 15),
		},
	}

	ai, err := AccruedInterest(bond, time.Now(), 0)
	test.AssertNoError(t, err)
	test.AssertEqual(t, 4.5, ai)
}

func Test_AccruedInterest_UsdBond(t *testing.T) {
	bond := Bond{
		Currency: "RUB",
		NominalCurrency: "USD",
	}
	bond.Coupons = []Coupon{
		{
			PerBondAmount:   45,
			CouponStartDate: time.Now(),
			CouponPeriod:    10,
			CouponEndDate:   time.Now().Add(time.Hour * 24 * 30),
			CouponDate:      time.Now().Add(time.Hour * 24 * 30),
		},
		{
			PerBondAmount:   70,
			CouponPeriod:    10,
			CouponStartDate: time.Now().Add(-time.Hour * 24),
			CouponEndDate:   time.Now().Add(-time.Hour * 12),
			CouponDate:      time.Now().Add(-time.Hour * 12),
		},
		{
			PerBondAmount:   25,
			CouponPeriod:    10,
			CouponStartDate: time.Now().Add(-time.Hour * 24 * 30),
			CouponEndDate:   time.Now().Add(-time.Hour * 24 * 15),
			CouponDate:      time.Now().Add(-time.Hour * 24 * 15),
		},
	}

	ai, err := AccruedInterest(bond, time.Now(), 80)
	test.AssertNoError(t, err)
	test.AssertEqual(t, 360, ai)
}
