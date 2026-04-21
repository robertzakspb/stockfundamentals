package bondportfolioanalysis

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	accountmvservice "github.com/compoundinvest/stockfundamentals/internal/application/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
)

func GeneratePortfolioOverview(filters []ydbfilter.YdbFilter) (string, error) {
	var sb strings.Builder

	sb.WriteString("Портфель по состоянию на ")
	sb.WriteString(timehelpers.TodayInDDMMYYYFormat())
	sb.WriteString("\n")
	sb.WriteString("\n")

	acountReturn, err := accountmvservice.GetAccountReturn(filters)
	if err != nil {
		return "", err
	}
	sb.WriteString("Текущая прибыль: ")
	if acountReturn.AbsoluteReturn >= 0 {
		sb.WriteString("+")
	} else {
		sb.WriteString("-")
	}
	sb.WriteString(fmt.Sprint(acountReturn.AbsoluteReturn))
	sb.WriteString("\n")

	//Adding the currency-based asset market values
	sb.WriteString("Стоимость активов: ")
	sb.WriteString("\n")

	mvs, err := accountmvservice.CalculateAccountMarketValue(acountReturn.AccountId, time.Now())
	if err != nil {
		return sb.String(), err
	}

	for _, mv := range mvs {
		sb.WriteString("  - ")
		sb.WriteString(forexservice.GetCurrencySymbol(mv.Currency))
		fmt.Fprint(&sb, mv.EodValue)
		sb.WriteString("\n")
	}

	//Adding the coupons payable withing the next seven days
	sb.WriteString("Выплачивамые на следующей неделе купоны: ")
	portfolio, err := bondportfolio.GetAllPositionLots()
	if err != nil {
		return sb.String(), err
	}
	portfolio, err = bondportfolio.PopulateLotsWithBonds(portfolio)
	if err != nil {
		return sb.String(), err
	}
	portfolio = bondportfolio.PopulateLotsWithCoupons(portfolio)

	for i := range portfolio {
		oneWeekFromNow := time.Now().AddDate(0, 0, 7)
		for _, coupon := range portfolio[i].Bond.Coupons {
			//Looking up coupons payable within the next seven days
			if timehelpers.DateIsLaterOrSameDate(coupon.CouponDate, time.Now()) && timehelpers.DateIsEarlierOrSameDate(coupon.CouponDate, oneWeekFromNow) {
				sb.WriteString("  - ")
				sb.WriteString(portfolio[i].Bond.Name)
				sb.WriteString(": ")
				sb.WriteString(forexservice.GetCurrencySymbol(portfolio[i].Bond.NominalCurrency))
				fmt.Fprint(&sb, coupon.PerBondAmount)
				sb.WriteString(" на ")
				sb.WriteString(strconv.Itoa(int(portfolio[i].Quantity)))
				sb.WriteString(" шт. = ")
				sb.WriteString(forexservice.GetCurrencySymbol(portfolio[i].Bond.NominalCurrency))
				fmt.Fprint(&sb, coupon.PerBondAmount*portfolio[i].Quantity)
			}
		}
	}

	return sb.String(), nil
}
