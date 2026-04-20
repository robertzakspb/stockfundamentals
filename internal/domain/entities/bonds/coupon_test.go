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

	_, err := AccruedInterest(bond, time.Now())

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

	_, err := AccruedInterest(bond, time.Now())

	test.AssertError(t, err)
}

func Test_AccruedInterest_OneDayBeforeNextCouponStartDate(t *testing.T) {
	bond := Bond{
		CouponCountPerYear: 12,
		NominalValue:       100,
		Currency:           "RUB",
		NominalCurrency:    "USD",
	}
	bond.Coupons = append(bond.Coupons, Coupon{
		CouponPeriod:    30,
		PerBondAmount:   0.62,
		CouponStartDate: time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC),
		CouponEndDate:   time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC),
		CouponDate:      time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC),
	})
	bond.Coupons = append(bond.Coupons, Coupon{
		CouponPeriod:  30,
		PerBondAmount: 0.62,

		CouponStartDate: time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC),
		CouponEndDate:   time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC),
		CouponDate:      time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC),
	})
	bond.Coupons = append(bond.Coupons, Coupon{
		CouponPeriod:    30,
		PerBondAmount:   0.62,
		CouponStartDate: time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC),
		CouponEndDate:   time.Date(2026, 5, 21, 0, 0, 0, 0, time.UTC),
		CouponDate:      time.Date(2026, 5, 21, 0, 0, 0, 0, time.UTC),
	})

	ai, err := AccruedInterest(bond, time.Now())

	test.AssertNoError(t, err)
	test.AssertEqual(t, 0, ai)
}

func Test_AccruedInterest_TodayIsCouponStartDate(t *testing.T) {
	bond := Bond{
		CouponCountPerYear: 12,
		NominalValue:       1000,
		Currency:           "RUB",
		NominalCurrency:    "RUB",
	}
	bond.Coupons = append(bond.Coupons, Coupon{
		CouponPeriod:    30,
		PerBondAmount:   10.89,
		CouponStartDate: time.Date(2026, 3, 21, 0, 0, 0, 0, time.UTC),
		CouponEndDate:   time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC),
		CouponDate:      time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC),
	})
	bond.Coupons = append(bond.Coupons, Coupon{
		CouponPeriod:    30,
		PerBondAmount:   10.89,
		CouponStartDate: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC),
		CouponEndDate:   time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
		CouponDate:      time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
	})

	ai, err := AccruedInterest(bond, time.Now())

	test.AssertNoError(t, err)
	test.AssertEqual(t, 0.36, ai)
}

func Test_AccruedInterest_TypicalScenario(t *testing.T) {
	bond := Bond{
		CouponCountPerYear: 12,
		NominalValue:       1000,
		Currency:           "RUB",
		NominalCurrency:    "USD",
	}
	bond.Coupons = append(bond.Coupons, Coupon{
		CouponPeriod:    30,
		PerBondAmount:   6.37,
		CouponStartDate: time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
		CouponEndDate:   time.Date(2026, 5, 17, 0, 0, 0, 0, time.UTC),
		CouponDate:      time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC),
	})

	ai, err := AccruedInterest(bond, time.Now())

	test.AssertNoError(t, err)
	test.AssertEqual(t, 0.85, ai)
}

func Test_AccruedInterest_OneDayBeforeEndDate(t *testing.T) {
	bond := Bond{
		CouponCountPerYear: 12,
		NominalValue:       1000,
		Currency:           "RUB",
		NominalCurrency:    "RUB",
	}
	bond.Coupons = append(bond.Coupons, Coupon{
		CouponPeriod:    30,
		PerBondAmount:   13.15,
		CouponStartDate: time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC),
		CouponEndDate:   time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
		CouponDate:      time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
	})

	ai, err := AccruedInterest(bond, time.Now())

	test.AssertNoError(t, err)
	test.AssertEqual(t, 12.71, ai)
}
