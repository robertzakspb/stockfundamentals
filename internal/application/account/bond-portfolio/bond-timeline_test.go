package bondportfolio

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_generateTimeLine_WithZeroLots(t *testing.T) {
	_, err := generateTimeLineForLots([]bonds.BondLot{})
	test.AssertError(t, err)
}

func Test_generateTimeLine_LotsWithNoBonds(t *testing.T) {
	lot := bonds.BondLot{
		Isin: "test",
	}
	_, err := generateTimeLineForLots([]bonds.BondLot{lot})
	test.AssertError(t, err)
}

func Test_generateTimeline_Positive(t *testing.T) {
	maturityDate := time.Date(2026, 12, 12, 0, 0, 0, 0, time.UTC)
	coupons := []bonds.Coupon{
		{
			PerBondAmount: 29.92,
			CouponDate:    time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			PerBondAmount: 30,
			CouponDate:    time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			PerBondAmount: 31,
			CouponDate:    time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			PerBondAmount: 32,
			CouponDate:    time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC),
		},
	}

	lot := bonds.BondLot{
		Isin:     "test",
		Quantity: 10,
		Bond: bonds.Bond{
			MaturityDate:    maturityDate,
			NominalValue:    1000,
			NominalCurrency: "RUB",
			Coupons:         coupons,
		},
	}

	timeline, err := generateTimeLineForLots([]bonds.BondLot{lot})
	test.AssertNoError(t, err)

	test.AssertEqual(t, 29.92, timeline[0].Amount)
	test.AssertEqual(t, 30, timeline[1].Amount)
	test.AssertEqual(t, 31, timeline[2].Amount)
	test.AssertEqual(t, 32, timeline[3].Amount)
	test.AssertEqual(t, time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC), timeline[0].Timestamp)
	test.AssertEqual(t, time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC), timeline[1].Timestamp)
	test.AssertEqual(t, time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC), timeline[2].Timestamp)
	test.AssertEqual(t, time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC), timeline[3].Timestamp)
	test.AssertEqual(t, "RUB", timeline[0].Currency)
	test.AssertEqual(t, "RUB", timeline[1].Currency)
	test.AssertEqual(t, "RUB", timeline[2].Currency)
	test.AssertEqual(t, "RUB", timeline[3].Currency)
	test.AssertEqual(t, "RUB", timeline[3].Currency)
	
	test.AssertEqual(t, lot.Bond.NominalValue*10, timeline[4].Amount)
	test.AssertEqual(t, lot.Bond.MaturityDate, timeline[4].Timestamp)
	test.AssertEqual(t, "RUB", timeline[4].Currency)
}
