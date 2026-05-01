package divcalendarservice

import (
	"time"

	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	divcal "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/dividend-calendar"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	stockportfolio "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/google/uuid"
)

func GetAccountDividendCalendar(accountIds uuid.UUIDs) (divcal.DividendCalendar, error) {
	accountIdFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "account_id",
		Condition:      ydbfilter.Contains,
		ConditionValue: ydbhelper.ConvertUUIDsToYdbList(accountIds),
	}
	portfolio, err := portfolio.GetAccountPortfolio([]ydbfilter.YdbFilter{accountIdFilter})
	if err != nil {
		return divcal.DividendCalendar{}, err
	}

	securityIds := stockportfolio.LotFigis(portfolio.Lots)
	securities, err := security_master.GetSecuritiesFilteredByFigi(securityIds)
	if err != nil {
		return divcal.DividendCalendar{}, err
	}

	filters := []ydbfilter.YdbFilter{{
		YqlColumnName:  "record_date",
		Condition:      ydbfilter.GreaterThanOrEqualTo,
		ConditionValue: ydbhelper.ConvertToYdbDate(time.Now()),
	}}

	figis := security_master.ExtractFigisFromSecurities(securities)
	filters = append(filters, ydbfilter.YdbFilter{
		YqlColumnName:  "stock_id",
		Condition:      ydbfilter.Contains,
		ConditionValue: ydbhelper.ConvertStringsToYdbList(figis),
	})

	relevantDivs, err := appdividend.GetFilteredDividends(filters)
	if err != nil {
		return divcal.DividendCalendar{}, err
	}

	relevantDivs = dividend.MatchDividendsWithStocks(relevantDivs, securities)

	divCal := divcal.DividendCalendar{AccountIds: accountIds}

	for _, relevantDiv := range relevantDivs {
		for _, lot := range portfolio.UniquePositions() {
			if relevantDiv.Figi == lot.Figi {
				divCal.FuturePayouts = append(divCal.FuturePayouts, dividend.Payout{
					Id:         uuid.New(),
					DividendId: relevantDiv.Id,
					AccountId:  lot.AccountId,
					Amount:     lot.Quantity * relevantDiv.ActualDPS,
					Dividend:   relevantDiv,
				})
			}
		}
	}

	return divCal, nil
}
