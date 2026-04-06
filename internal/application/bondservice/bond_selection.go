package bondservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func GetRussianGovernmentBondsWithFixedOrConstantCoupon() ([]bonds.Bond, error) {
	governmentFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "name",
		Condition:      ydbfilter.Like,
		ConditionValue: types.TextValue("%ОФЗ%"),
	}
	amortizationFilter :=
		ydbfilter.YdbFilter{
			YqlColumnName:  "has_amortization",
			Condition:      ydbfilter.Equal,
			ConditionValue: types.BoolValue(false),
		}

	bondList, err := GetFilteredBonds([]ydbfilter.YdbFilter{governmentFilter, amortizationFilter})
	if err != nil {
		return bondList, err
	}

	bondList = PopulateBondCoupons(bondList)

	bondList = GetOnlyBondsWithFixedOrConstantCoupons(bondList)

	bondsWithYtm := CalculateYtmForBonds(bondList)
	return bondsWithYtm, nil
}

func GetQuasiForeignBonds() ([]bonds.Bond, error) {
	foreignNominalFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "nominal_currency",
		Condition:      ydbfilter.NotEqual,
		ConditionValue: types.TextValue("RUB"),
	}
	rubleCurrencyFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "currency",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue("RUB"),
	}

	bondList, err := GetFilteredBonds([]ydbfilter.YdbFilter{foreignNominalFilter, rubleCurrencyFilter})
	if err != nil {
		return bondList, err
	}

	bondList = PopulateBondsWithCouponsAndCalculateYtm(bondList)

	return bondList, nil
}
