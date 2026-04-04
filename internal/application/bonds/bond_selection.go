package bondservice

import (
	"github.com/compoundinvest/invest-core/quote/bondquote"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
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

	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log("Failed to initialize the configuration file", logger.ALERT)
		return []bonds.Bond{}, nil
	}

	quotes, err := bondquote.FetchQuotesForFigis(GetBondFigis(bondList), config)

	bondsWithYtm := CalculateYtmForBonds(bondList, quotes)
	return bondsWithYtm, nil
}

/*
SELECT `id`, `figi`, `isin`, `lot`, `currency`, `name`, `country_of_risk`, `real_exchange`, `coupon_count_per_year`, `nominal_value`, `nominal_currency`, `initial_nominal_value`, `initial_nominal_currency`, `placement_price`, `placement_currency`, `accumulated_coupon_income`, `issue_size`, `issue_size_plan`, `has_floating_coupon`, `is_perpetual`, `has_amortization`, `is_available_for_iis`, `is_for_qualified_investors`, `is_subordinated`, `risk_level`, `bond_type`, `call_option_exercise_date`, `registration_date`, `placement_date`, `maturity_date`
FROM `bonds/bond`
WHERE name = 'ОФЗ 26250'
ORDER BY maturity_date DESC
*/
