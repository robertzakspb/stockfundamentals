package bondservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func GetRussianGovernmentBondsWithFixedCoupon() ([]bonds.Bond, error) {
	governmentFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "name",
		Condition:      ydbfilter.Like,
		ConditionValue: types.TextValue("%ОФЗ%"),
	}

	bonds, err := GetFilteredBonds([]ydbfilter.YdbFilter{governmentFilter})
	if err != nil {
		return bonds, err
	}

	return bonds, nil
}

/*
SELECT `id`, `figi`, `isin`, `lot`, `currency`, `name`, `country_of_risk`, `real_exchange`, `coupon_count_per_year`, `nominal_value`, `nominal_currency`, `initial_nominal_value`, `initial_nominal_currency`, `placement_price`, `placement_currency`, `accumulated_coupon_income`, `issue_size`, `issue_size_plan`, `has_floating_coupon`, `is_perpetual`, `has_amortization`, `is_available_for_iis`, `is_for_qualified_investors`, `is_subordinated`, `risk_level`, `bond_type`, `call_option_exercise_date`, `registration_date`, `placement_date`, `maturity_date`
FROM `bonds/bond`
WHERE name = 'ОФЗ 26250'
ORDER BY maturity_date DESC
*/
