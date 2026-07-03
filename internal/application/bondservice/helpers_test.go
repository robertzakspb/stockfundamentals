package bondservice

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_GetBondFigis(t *testing.T) {
	bondList := []bonds.Bond{
		{
			Figi: "figi1",
		},
		{
			Figi: "figi2",
		},
		{
			Figi: "figi3",
		},
		{
			Figi: "",
		},
	}

	figis := GetBondFigis(&bondList)

	test.AssertEqual(t, 3, len(figis))
	test.AssertEqual(t, "figi1", figis[0])
	test.AssertEqual(t, "figi2", figis[1])
	test.AssertEqual(t, "figi3", figis[2])
}

func Test_GetOnlyBondsWithFixedOrConstantCoupons(t *testing.T) {
	bondList := []bonds.Bond{
		{
			Figi: "figi1",
			Coupons: []bonds.Coupon{
				{
					CouponType: bonds.CouponType_COUPON_TYPE_CONSTANT,
				},
			},
		},
		{
			Figi: "figi1",
			Coupons: []bonds.Coupon{
				{
					CouponType: bonds.CouponType_COUPON_TYPE_CONSTANT,
				},
			},
		},
		{
			Figi: "figi1",
			Coupons: []bonds.Coupon{
				{
					CouponType: bonds.CouponType_COUPON_TYPE_DISCOUNT,
				},
			},
		},
	}

	bondList = GetOnlyBondsWithFixedOrConstantCoupons(bondList)

	test.AssertEqual(t, 2, len(bondList))
	test.AssertEqual(t, bonds.CouponType_COUPON_TYPE_CONSTANT, bondList[0].Coupons[0].CouponType)
	test.AssertEqual(t, bonds.CouponType_COUPON_TYPE_CONSTANT, bondList[1].Coupons[0].CouponType)
}

func Test_AllCurrencyPairsInBondList(t *testing.T) {
	bondList := []bonds.Bond{
		{
			Figi:            "figi1",
			NominalCurrency: "USD",
			Currency:        "RUB",
		},
		{
			Figi:            "figi2",
			NominalCurrency: "EUR",
			Currency:        "RUB",
		},
		{
			Figi:            "figi3",
			NominalCurrency: "RUB",
			Currency:        "RUB",
		},
		{
			Figi:            "",
			NominalCurrency: "USD",
			Currency:        "RUB",
		},
	}

	currencyPairs := AllCurrencyPairsInBondList(bondList)

	test.AssertEqual(t, 2, len(currencyPairs))
	test.AssertEqual(t, "USD/RUB", currencyPairs[0])
	test.AssertEqual(t, "EUR/RUB", currencyPairs[1])
}

func Test_MatchCouponsWithBonds(t *testing.T) {
	bondList := []bonds.Bond{
		{
			Figi: "figi1",
		},
		{
			Figi: "figi2",
		},
		{
			Figi: "figi3",
		},
		{
			Figi: "",
		},
	}
	couponList := []bonds.Coupon{
		{
			Figi:         "figi3",
			CouponNumber: 1,
		},
		{
			Figi:         "figi3",
			CouponNumber: 2,
		},
		{
			Figi:         "figi2",
			CouponNumber: 1,
		},
		{
			Figi:         "figi1",
			CouponNumber: 2,
		},
		{
			Figi:         "figi1",
			CouponNumber: 3,
		},
		{
			Figi:         "",
			CouponNumber: 2,
		},
	}

	bondList = MatchCouponsWithBonds(couponList, bondList)

	test.AssertEqual(t, 4, len(bondList))

	test.AssertEqual(t, "figi1", bondList[0].Figi)
	test.AssertEqual(t, 2, len(bondList[0].Coupons))
	test.AssertEqual(t, "figi1", bondList[0].Coupons[0].Figi)
	test.AssertEqual(t, 2, bondList[0].Coupons[0].CouponNumber)
	test.AssertEqual(t, "figi1", bondList[0].Coupons[1].Figi)
	test.AssertEqual(t, 3, bondList[0].Coupons[1].CouponNumber)

	test.AssertEqual(t, "figi2", bondList[1].Figi)
	test.AssertEqual(t, 1, len(bondList[1].Coupons))
	test.AssertEqual(t, "figi2", bondList[1].Coupons[0].Figi)
	test.AssertEqual(t, 1, bondList[1].Coupons[0].CouponNumber)

	test.AssertEqual(t, "figi3", bondList[2].Figi)
	test.AssertEqual(t, 2, len(bondList[2].Coupons))
	test.AssertEqual(t, "figi3", bondList[2].Coupons[0].Figi)
	test.AssertEqual(t, "figi3", bondList[2].Coupons[1].Figi)
	test.AssertEqual(t, 1, bondList[2].Coupons[0].CouponNumber)
	test.AssertEqual(t, 2, bondList[2].Coupons[1].CouponNumber)

	test.AssertEqual(t, "", bondList[3].Figi)
	test.AssertEqual(t, 2, len(bondList[2].Coupons))
}
