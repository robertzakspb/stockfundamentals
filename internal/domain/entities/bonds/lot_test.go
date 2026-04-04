package bonds

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_PricePerUnitPercentage_Negative(t *testing.T) {
	lot := BondLot{}

	pricePercentage := lot.PricePerUnitPercentage(0)

	test.AssertEqual(t, -1, pricePercentage)
}

func Test_PricePerUnitPercentage_NonRubPrice(t *testing.T) {
	lot := BondLot{
		PricePerUnit: 857,
	}

	pricePercentage := lot.PricePerUnitPercentage(1000)

	test.AssertEqual(t, 0.857, pricePercentage)
}

func Test_PricePerUnitPercentage_RubPrice(t *testing.T) {
	lot := BondLot{
		PricePerUnitInRUB: 80,
	}

	pricePercentage := lot.PricePerUnitPercentage(100)

	test.AssertEqual(t, 0.8, pricePercentage)
}

func Test_CouponPayoutForPosition(t *testing.T) {
	lot := BondLot{
		Quantity: 10,
	}

	payout := lot.CouponPayoutForPosition(Coupon{PerBondAmount: 25})

	test.AssertEqual(t, 250, payout)
}

func Test_TotalPrincipalRedemption(t *testing.T) {
	lot := BondLot{
		Quantity: 10,
	}

	totalRedemption := lot.TotalPrincipalRedemption(Bond{NominalValue: 1005})

	test.AssertEqual(t, 10050, totalRedemption)
}
