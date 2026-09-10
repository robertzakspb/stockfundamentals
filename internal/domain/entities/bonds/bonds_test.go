package bonds

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
	"github.com/google/uuid"
)

// Helper
func generateMockValidBond() Bond {
	bond := Bond{
		Id:                      uuid.New(),
		Figi:                    "testFigi",
		Isin:                    "testIsin",
		Lot:                     10,
		Currency:                "USD",
		CouponCountPerYear:      10,
		MaturityDate:            time.Now(),
		NominalValue:            1000,
		NominalCurrency:         "EUR",
		InitialNominalValue:     1000,
		InitialNominalCurrency:  "EUR",
		PlacementPrice:          1005,
		PlacementCurrency:       "EUR",
		AccruedInterest:         10,
		IssueSize:               1_000_000,
		IssueSizePlan:           5_000_000,
		RiskLevel:               HIGH_RISK_LEVEL,
		BondType:                BondType_BOND_TYPE_UNSPECIFIED,
		CallOptionExerciseDate:  time.Now(),
		SimpleYieldToMaturity:   14.3,
		SimpleYieldToCallOption: 8.3,
		QuoteInPercentage:       90.3,
		Coupons: []Coupon{
			{
				Figi:            "testFigi",
				CouponDate:      time.Now(),
				RecordDate:      time.Now().AddDate(0, 0, -1),
				PerBondAmount:   10.3,
				CouponType:      CouponType_COUPON_TYPE_FIX,
				CouponStartDate: time.Now().AddDate(0, 0, -5),
				CouponEndDate:   time.Now(),
				CouponPeriod:    30,
			},
		},
	}

	return bond
}

func Test_validate_ProperBond(t *testing.T) {
	validBond := generateMockValidBond()

	err := validBond.Validate()

	test.AssertNoError(t, err)
}

func Test_Validate_NilId(t *testing.T) {
	bond := generateMockValidBond()
	bond.Id = uuid.Nil

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_MissingFigi(t *testing.T) {
	bond := generateMockValidBond()
	bond.Figi = ""

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_MissingIsin(t *testing.T) {
	bond := generateMockValidBond()
	bond.Isin = ""

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidLot(t *testing.T) {
	bond := generateMockValidBond()
	bond.Lot = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidCurrency(t *testing.T) {
	bond := generateMockValidBond()
	bond.Currency = "ABC"

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidCouponCount(t *testing.T) {
	bond := generateMockValidBond()
	bond.CouponCountPerYear = -1

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidMaturityDate(t *testing.T) {
	bond := generateMockValidBond()
	bond.MaturityDate = time.Time{}

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidNominalValue(t *testing.T) {
	bond := generateMockValidBond()
	bond.NominalValue = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidNominalCurrency(t *testing.T) {
	bond := generateMockValidBond()
	bond.NominalCurrency = "WER"

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidInitialNominalValue(t *testing.T) {
	bond := generateMockValidBond()
	bond.InitialNominalValue = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidInitialNominalCurrency(t *testing.T) {
	bond := generateMockValidBond()
	bond.InitialNominalCurrency = "TEST"

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidPlacementPrice(t *testing.T) {
	bond := generateMockValidBond()
	bond.PlacementPrice = 0.0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidPlacementCurrency(t *testing.T) {
	bond := generateMockValidBond()
	bond.PlacementCurrency = "PPP"

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidAccruedInterest(t *testing.T) {
	bond := generateMockValidBond()
	bond.AccruedInterest = -1

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidIssueSize(t *testing.T) {
	bond := generateMockValidBond()
	bond.IssueSize = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_Validate_InvalidIssueSizePlan(t *testing.T) {
	bond := generateMockValidBond()
	bond.IssueSizePlan = 0

	err := bond.Validate()

	test.AssertError(t, err)
}

func Test_HasCallOption_Negative(t *testing.T) {
	bond := generateMockValidBond()
	bond.CallOptionExerciseDate = time.Time{}

	test.AssertFalse(t, bond.HasCallOption())
}

func Test_HasCallOption_Positive(t *testing.T) {
	bond := generateMockValidBond()
	bond.CallOptionExerciseDate = time.Now()

	test.AssertTrue(t, bond.HasCallOption())
}

func Test_IsRubleBond_Negative(t *testing.T) {
	bond := generateMockValidBond()
	bond.Currency = "RUB"
	bond.NominalCurrency = "USD"

	test.AssertFalse(t, bond.IsRubleBond())
}

func Test_IsRubleBond_Positive(t *testing.T) {
	bond := generateMockValidBond()
	bond.Currency = "RUB"
	bond.NominalCurrency = "RUB"

	test.AssertTrue(t, bond.IsRubleBond())

}

func Test_CurrentCouponYield_Positive(t *testing.T) {
	bond := generateMockValidBond()

	test.AssertEqual(t, 0.11406423034330011, bond.CurrentCouponYield())
}

func Test_MarketPriceInCurrency_Positive(t *testing.T) {
	bond := generateMockValidBond()

	test.AssertEqual(t, 903, bond.MarketPriceInCurrency(90.3))
}

func Test_HasFixedCoupon_Negative_NoCoupons(t *testing.T) {
	b := Bond{}

	_, err := b.HasFixedCoupon()

	test.AssertError(t, err)
}

func Test_HasFixedCoupon_Negative_FloatingCoupon(t *testing.T) {
	b := Bond{}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType: CouponType_COUPON_TYPE_FLOATING,
	})

	isFixed, err := b.HasFixedCoupon()

	test.AssertNoError(t, err)
	test.AssertFalse(t, isFixed)
}

func Test_HasFixedCoupon_Positive_FixedCoupon(t *testing.T) {
	b := Bond{}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType: CouponType_COUPON_TYPE_FIX,
	})

	isFixed, err := b.HasFixedCoupon()

	test.AssertNoError(t, err)
	test.AssertTrue(t, isFixed)
}

func Test_HasFixedCoupon_Positive_ConstantCoupon(t *testing.T) {
	b := Bond{}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType: CouponType_COUPON_TYPE_CONSTANT,
	})

	isFixed, err := b.HasFixedCoupon()

	test.AssertNoError(t, err)
	test.AssertTrue(t, isFixed)
}

func Test_CurrentCouponYield_Negative_NoCoupons(t *testing.T) {
	b := Bond{
		QuoteInPercentage: 98,
		NominalValue:      1000,
	}

	yield := b.CurrentCouponYield()

	test.AssertEqual(t, 0, yield)
}

func Test_CurrentCouponYield_Negative_NoQuote(t *testing.T) {
	b := Bond{}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType: CouponType_COUPON_TYPE_CONSTANT,
	})

	yield := b.CurrentCouponYield()

	test.AssertEqual(t, 0, yield)
}

func Test_CurrentCouponYield_Negative_FloatingCoupon(t *testing.T) {
	b := Bond{}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType: CouponType_COUPON_TYPE_FLOATING,
	})

	yield := b.CurrentCouponYield()

	test.AssertEqual(t, 0, yield)
}

func Test_TotalFutureCoupons_Negative_MissingCoupons(t *testing.T) {
	b := Bond{}

	_, err := b.TotalFutureCoupons(true)

	test.AssertError(t, err)
}

func Test_TotalFutureCoupons_Negative_FloatingCoupon(t *testing.T) {
	b := Bond{}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType: CouponType_COUPON_TYPE_FLOATING,
	})

	_, err := b.TotalFutureCoupons(true)

	test.AssertError(t, err)
}

func Test_TotalFutureCoupons_Positive_UntillCallOption(t *testing.T) {
	b := Bond{
		CallOptionExerciseDate: time.Date(2100, 11, 10, 0, 0, 0, 0, time.UTC),
	}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 9, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 10, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})

	tci, err := b.TotalFutureCoupons(true)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 40, tci)
}

func Test_TotalFutureCoupons_Positive_UntillMaturity(t *testing.T) {
	b := Bond{
		MaturityDate: time.Date(2100, 11, 10, 0, 0, 0, 0, time.UTC),
	}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 9, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 10, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})

	tci, err := b.TotalFutureCoupons(true)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 60, tci)
}

func Test_TotalFutureCashflows(t *testing.T) {
	b := Bond{
		MaturityDate:      time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC),
		AccruedInterest:   7,
		QuoteInPercentage: 98,
		NominalValue:      1000,
	}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 9, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 10, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	expectedCashflow := -b.MarketPriceInCurrency(b.QuoteInPercentage) - 7 + 20 + 20 + 20 + 1000

	actualCashflow, err := b.TotalFutureCashflows(false, 980)

	test.AssertNoError(t, err)
	test.AssertEqual(t, expectedCashflow, actualCashflow)
}

func Test_FutureCouponsTillDate_Negative_NoCoupons(t *testing.T) {
	b := Bond{}

	_, err := b.FutureCouponsTillDate(time.Now())

	test.AssertError(t, err)

}

func Test_FutureCouponsTillDate_Positive_OneCouponIsPaidToday(t *testing.T) {
	b := Bond{
		MaturityDate:      time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC),
		AccruedInterest:   7,
		QuoteInPercentage: 98,
		NominalValue:      1000,
	}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Now(),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 10, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})

	coupons, err := b.FutureCouponsTillDate(time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC))

	test.AssertNoError(t, err)
	test.AssertEqual(t, 2, len(coupons))
	test.AssertEqual(t, time.Date(2100, 10, 16, 0, 0, 0, 0, time.UTC), coupons[0].CouponDate)
	test.AssertEqual(t, time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC), coupons[1].CouponDate)
}

func Test_FutureCouponsTillDate_Positive_NoFutureCoupons(t *testing.T) {
	b := Bond{
		MaturityDate:      time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC),
		AccruedInterest:   7,
		QuoteInPercentage: 98,
		NominalValue:      1000,
	}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Now(),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Now(),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Now(),
		PerBondAmount: 20,
	})

	coupons, err := b.FutureCouponsTillDate(time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC))

	test.AssertNoError(t, err)
	test.AssertEqual(t, 0, len(coupons))
}

func Test_FutureCashflowsWithDates_Negative_MissingCoupons(t *testing.T) {
	b := Bond{}

	_, _, err := b.FutureCashflowsWithDates(false, 100, time.Now())

	test.AssertError(t, err)
}

func Test_FutureCashflowsWithDates_Positive(t *testing.T) {
	b := Bond{
		MaturityDate:      time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC),
		AccruedInterest:   7,
		QuoteInPercentage: 98,
		NominalValue:      1000,
	}
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 9, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 10, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	b.Coupons = append(b.Coupons, Coupon{
		CouponType:    CouponType_COUPON_TYPE_FIX,
		CouponDate:    time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC),
		PerBondAmount: 20,
	})
	// expectedCashflow := -b.MarketPriceInCurrency(b.QuoteInPercentage) - 7 + 20 + 20 + 20 + 1000

	cashflows, dates, err := b.FutureCashflowsWithDates(false, 980, time.Now())

	test.AssertNoError(t, err)
	test.AssertEqual(t, 5, len(cashflows))
	test.AssertEqual(t, 5, len(dates))
	test.AssertEqual(t, -b.MarketPriceInCurrency(b.QuoteInPercentage)-7, cashflows[0])
	test.AssertEqual(t, 20, cashflows[1])
	test.AssertEqual(t, 20, cashflows[2])
	test.AssertEqual(t, 20, cashflows[3])
	test.AssertEqual(t, 1000, cashflows[4])
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now(), dates[0]))
	test.AssertEqual(t, time.Date(2100, 9, 16, 0, 0, 0, 0, time.UTC), dates[1])
	test.AssertEqual(t, time.Date(2100, 10, 16, 0, 0, 0, 0, time.UTC), dates[2])
	test.AssertEqual(t, time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC), dates[3])
	test.AssertEqual(t, time.Date(2100, 11, 16, 0, 0, 0, 0, time.UTC), dates[4])

}
