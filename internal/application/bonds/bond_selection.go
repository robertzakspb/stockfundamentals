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

	bondList, err := GetFilteredBonds([]ydbfilter.YdbFilter{governmentFilter})
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

	bondList = PopulateBondCoupons(bondList)

	bondList = CalculateYtmForBonds(bondList)

	return bondList, nil
}
