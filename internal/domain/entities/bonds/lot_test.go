package bonds

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

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

func Test_CurrentProfitOrLossPercentage_Positive_Profit(t *testing.T) {
	lot := BondLot{
		PricePerUnitPercentage: 92.4,
		Bond: Bond{
			QuoteInPercentage: 100,
		},
	}

	currentPL := lot.CurrentProfitOrLossPercentage()

	test.AssertEqualFloat(t, 0.08225, currentPL, 0.0001)
}

func Test_CurrentProfitOrLossPercentage_Positive_Loss(t *testing.T) {
	lot := BondLot{
		PricePerUnitPercentage: 100,
		Bond: Bond{
			QuoteInPercentage: 95,
		},
	}

	currentPL := lot.CurrentProfitOrLossPercentage()

	test.AssertEqualFloat(t, -0.05, currentPL, 0.0001)
}

func Test_CurrentProfitOrLossPercentage_Negative_MissingCost(t *testing.T) {
	lot := BondLot{
		PricePerUnitPercentage: 0,
		Bond: Bond{
			QuoteInPercentage: 95,
		},
	}

	currentPL := lot.CurrentProfitOrLossPercentage()

	test.AssertEqualFloat(t, 0, currentPL, 0.0001)
}
