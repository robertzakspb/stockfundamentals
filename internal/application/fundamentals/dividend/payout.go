package appdividend

import (
	"time"

	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/google/uuid"
)

func GetDividendPayoutsForAccount(accountId uuid.UUID) ([]dividend.Payout, error) {
	lots, err := portfolio.GetAccountPortfolio(accountId)
	if err != nil {
		return []dividend.Payout{}, err
	}
	figis := stockportfolio.LotFigis(lots.Lots)

	filters := []ydbfilter.YdbFilter{
		{
			YqlColumnName:  "record_date",
			Condition:      ydbfilter.GreaterThan,
			ConditionValue: ydbhelper.ConvertToOptionalYDBdate(time.Now().AddDate(0, 0, -21)), //We fetch dividends with the record date as early as 3 weeks ago to account for the payout delay (~2 weeks on average)
		},
		{
			YqlColumnName:  "stock_id",
			Condition:      ydbfilter.Contains,
			ConditionValue: ydbhelper.ConvertStringsToYdbList(figis),
		},
	}

	dividends, err := GetFilteredDividends(filters)
	if err != nil {
		return []dividend.Payout{}, err
	}
	dividends, err = PopulateDividendStocks(dividends)

	payouts := dividend.MakePayoutsFromDividendsAndLots(lots.Lots, dividends)
	if err != nil {
		return []dividend.Payout{}, err
	}

	return payouts, nil
}
