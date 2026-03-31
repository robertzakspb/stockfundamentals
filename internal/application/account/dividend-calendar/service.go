package divcalendarservice

import (
	"time"

	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	divcal "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/dividend-calendar"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
)

func GetAccountDividendCalendar(accountIds uuid.UUIDs) (divcal.DividendCalendar, error) {
	portfolio, err := portfolio.GetAccountPortfolio(accountIds)
	if err != nil {
		return divcal.DividendCalendar{}, err
	}

	securityIds := portfolio.Securities()
	securities, err := security_master.GetSecuritiesById(securityIds)
	if err != nil {
		return divcal.DividendCalendar{}, err
	}

	filters := []ydbfilter.YdbFilter{{
		YqlColumnName:  "record_date",
		Condition:      ydbfilter.GreaterThanOrEqualTo,
		ConditionValue: shared.ConvertToYdbDate(time.Now()),
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

	divCal := divcal.DividendCalendar{AccountIds: accountIds}

	for _, relevantDiv := range relevantDivs {
		for _, lot := range portfolio.Lots {
			securityId, err := uuid.Parse(lot.SecurityId)
			if err != nil {
				logger.Log("Unexpectedly failed to parse a UUID from security ID: "+lot.SecurityId, logger.ERROR)
				continue
			}
			if relevantDiv.Id == securityId {
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
