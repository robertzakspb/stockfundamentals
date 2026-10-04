package appdividend

import (
	"sync"
	"time"

	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	"github.com/compoundinvest/stockfundamentals/internal/application/account/transactionprocessor"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
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

	var wg sync.WaitGroup
	transactions := []transaction.Transaction{}
	wg.Go(func() {
		dividends, err = PopulateDividendStocks(dividends)
	})
	wg.Go(func() {
		filters := []ydbfilter.YdbFilter{
			{
				YqlColumnName:  "timestamp",
				Condition:      ydbfilter.GreaterThan,
				ConditionValue: ydbhelper.ConvertToOptionalYDBdate(time.Now().AddDate(0, 0, -30)),
			},
			{
				YqlColumnName:  "type",
				Condition:      ydbfilter.Equal,
				ConditionValue: types.TextValue("DIVIDEND"),
			},
		}
		transactions, err = transactionprocessor.GetFilteredTransactions(shared.ParsedApiQuery{Filters: filters})
	})
	wg.Wait()

	payouts := dividend.MakePayoutsFromDividendsAndLots(lots.Lots, dividends, transactions)
	if err != nil {
		return []dividend.Payout{}, err
	}

	return payouts, nil
}
