package bondsapi

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
	"github.com/google/uuid"
)

func Test_mapBondsToDTOs(t *testing.T) {
	bonds := []bonds.Bond{
		{
			Id:                     uuid.New(),
			Figi:                   "testFigi",
			Isin:                   "testIsin",
			Lot:                    10,
			Currency:               "USD",
			CouponCountPerYear:     10,
			MaturityDate:           time.Now(),
			NominalValue:           1000,
			NominalCurrency:        "EUR",
			InitialNominalValue:    5000,
			InitialNominalCurrency: "USD",
			PlacementPrice:         1005,
			PlacementCurrency:      "EUR",
			AccruedInterest:        5.7,
			IssueSize:              1_000_000,
			IssueSizePlan:          5_000_000,
			RiskLevel:              bonds.HIGH_RISK_LEVEL,
			BondType:               bonds.BondType_BOND_TYPE_UNSPECIFIED,
			CallOptionExerciseDate: time.Now(),
			YieldToMaturity:        14.3,
			YieldToCallOption:      8.3,
			Coupons: []bonds.Coupon{
				{
					Figi:            "testFigi",
					CouponDate:      time.Now(),
					RecordDate:      time.Now().AddDate(0, 0, -1),
					PerBondAmount:   10.3,
					CouponType:      bonds.CouponType_COUPON_TYPE_FIX,
					CouponStartDate: time.Now().AddDate(0, 0, -5),
					CouponEndDate:   time.Now(),
					CouponPeriod:    30,
				},
			},
		},
	}

	mappedDtos := mapBondsToDTOs(bonds, true)

	test.AssertEqual(t, 1, len(mappedDtos))

	test.AssertEqual(t, "testFigi", mappedDtos[0].Figi)
	test.AssertEqual(t, "testIsin", mappedDtos[0].Isin)
	test.AssertEqual(t, 10, mappedDtos[0].Lot)
	test.AssertEqual(t, "USD", mappedDtos[0].Currency)
	test.AssertEqual(t, 10, mappedDtos[0].CouponCountPerYear)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now(), mappedDtos[0].MaturityDate))
	test.AssertEqual(t, 1000, mappedDtos[0].NominalValue)
	test.AssertEqual(t, "EUR", mappedDtos[0].NominalCurrency)
	test.AssertEqual(t, 5000, mappedDtos[0].InitialNominalValue)
	test.AssertEqual(t, 1005, mappedDtos[0].PlacementPrice)
	test.AssertEqual(t, "EUR", mappedDtos[0].PlacementCurrency)
	test.AssertEqual(t, 5.7, mappedDtos[0].AccruedInterest)
	test.AssertEqual(t, 1_000_000, mappedDtos[0].IssueSize)
	test.AssertEqual(t, 5_000_000, mappedDtos[0].IssueSizePlan)
	test.AssertEqual(t, "HIGH_RISK_LEVEL", mappedDtos[0].RiskLevel)
	test.AssertEqual(t, "BOND_TYPE_UNSPECIFIED", mappedDtos[0].BondType)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now(), mappedDtos[0].CallOptionExerciseDate))
	test.AssertEqual(t, 14.3, mappedDtos[0].YieldToMaturity)
	test.AssertEqual(t, 8.3, mappedDtos[0].YieldToCallOption)

	test.AssertEqual(t, 1, len(mappedDtos[0].Coupons))

	test.AssertEqual(t, "testFigi", mappedDtos[0].Coupons[0].Figi)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now(), mappedDtos[0].Coupons[0].Date))
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -1), mappedDtos[0].Coupons[0].RecordDate))
	test.AssertEqual(t, 10.3, mappedDtos[0].Coupons[0].PerBondAmount)
	test.AssertEqual(t, "COUPON_TYPE_FIX", mappedDtos[0].Coupons[0].CouponType)
	test.AssertEqual(t, 30, mappedDtos[0].Coupons[0].CouponPeriodInDays)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -5), mappedDtos[0].Coupons[0].CouponStartDate))
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now(), mappedDtos[0].Coupons[0].CouponEndDate))
}
