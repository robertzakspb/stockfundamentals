package bondservice

import (
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

func GetRussianGovernmentBondsWithFixedOrConstantCoupon() ([]bonds.Bond, error) {
	governmentFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "name",
		Condition:      ydbfilter.Like,
		ConditionValue: types.TextValue("%ОФЗ%"),
	}
	currencyFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "nominal_currency",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue("RUB"),
	}
	//Removing amortized bonds because YTM cannot be calculated for them
	amortizationFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "has_amortization",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.BoolValue(false),
	}
	//Removing bonds with indexed nominal value because YTM cannot be calculated for them
	nominalValueFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "nominal_value",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.DoubleValue(1000),
	}

	bondList, err := security_master.GetFilteredBonds([]ydbfilter.YdbFilter{governmentFilter, amortizationFilter, nominalValueFilter, currencyFilter})
	if err != nil {
		return bondList, err
	}

	bondList = PopulateBondsWithCouponsAndCalculateYtm(bondList)

	bondList = GetOnlyBondsWithFixedOrConstantCoupons(bondList)

	return bondList, nil
}

func GetQuasiForeignBonds() ([]bonds.Bond, error) {
	foreignNominalFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "nominal_currency",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue("USD"),
	}
	rubleCurrencyFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "currency",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue("RUB"),
	}
	riskFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "risk_level",
		Condition:      ydbfilter.NotEqual,
		ConditionValue: types.TextValue("HIGH_RISK_LEVEL"),
	}
	countryFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "country_of_risk",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.TextValue("RU"),
	}

	bondList, err := security_master.GetFilteredBonds([]ydbfilter.YdbFilter{foreignNominalFilter, rubleCurrencyFilter, riskFilter, countryFilter})
	if err != nil {
		return bondList, err
	}

	bondList = PopulateBondsWithCouponsAndCalculateYtm(bondList)

	return bondList, nil
}
