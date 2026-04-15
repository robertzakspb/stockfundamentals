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
