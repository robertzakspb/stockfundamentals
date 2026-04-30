package bondportfolioanalysis

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	accountmvservice "github.com/compoundinvest/stockfundamentals/internal/application/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/compoundinterest"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	stringhelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/string-helpers"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
	"github.com/google/uuid"
)

func GeneratePortfolioOverview(filters []ydbfilter.YdbFilter) (string, error) {
	var sb strings.Builder

	//Adding the currency-based asset market values
	sb.WriteString("Стоимость активов: ")

	accountReturn, err := accountmvservice.GetAccountReturn(filters, "RUB")
	if err != nil {
		return sb.String(), err
	}
	mv, _ := stringhelpers.BeautifyNumber(accountReturn.EndDateMV)
	sb.WriteString(mv)

	sb.WriteString("\n")

	sb.WriteString("Разбиение по активам: ")
	sb.WriteString("\n")
	hardcodedAccountId := uuid.MustParse("129274f9-ee80-4e74-aa1c-fea578bac6e6")
	mvs, err := accountmvservice.CalculateAccountMarketValue(hardcodedAccountId, time.Now())
	if err != nil {
		return sb.String(), err
	}

	for _, mv := range mvs {
		sb.WriteString("  - ")
		sb.WriteString(forexservice.GetCurrencySymbol(mv.Currency))
		beautifiedEOD, _ := stringhelpers.BeautifyNumber(mv.EodValue)
		fmt.Fprint(&sb, beautifiedEOD)
		sb.WriteString("\n")
	}

	for _, mv := range mvs {
		//Adding the current profit in the required currencies
		accountReturn, err := accountmvservice.GetAccountReturn(filters, mv.Currency)
		if err != nil {
			return "", err
		}
		generateAccountReturnOverview(&sb, accountReturn)
	}

	err = addNextWeekCoupons(&sb)
	if err != nil {
		return sb.String(), err
	}

	writeAnalysisToFile(&sb)

	return sb.String(), nil
}

func generateAccountReturnOverview(sb *strings.Builder, accountReturn accountmvdomain.Return) *strings.Builder {

	sb.WriteString("Текущая прибыль на ")
	sb.WriteString(timehelpers.TodayInDDMMYYYFormat())
	sb.WriteString(": ")
	sb.WriteString("\n")

	sb.WriteString("  - В ")
	sb.WriteString(forexservice.GetCurrencySymbol(accountReturn.Currency))
	sb.WriteString(": ")

	absoluteReturn, _ := stringhelpers.BeautifyNumber(accountReturn.AbsoluteReturn)
	sb.WriteString(absoluteReturn)

	sb.WriteString("( ")
	absoluteReturnPercentage, _ := stringhelpers.BeatufityPercentage(accountReturn.AbsoluteReturnPercentage)
	sb.WriteString(absoluteReturnPercentage)
	if accountReturn.AbsoluteReturn <= 0 {
		sb.WriteString(")")
		return sb
	}

	sb.WriteString("; или ")
	annualized := compoundinterest.CalcAnnualizedReturn(accountReturn.AbsoluteReturnPercentage, accountReturn.StartDate, accountReturn.EndDate)
	annualizedFormatted, _ := stringhelpers.BeatufityPercentage(annualized)
	sb.WriteString(annualizedFormatted)
	sb.WriteString(" годовых)")

	sb.WriteString("\n")

	return sb
}

func addNextWeekCoupons(sb *strings.Builder) error {
	//Adding the coupons payable withing the next seven days
	const header = "Выплачивамые на следующей неделе купоны: "
	sb.WriteString(header)
	portfolio, err := bondportfolio.GetAllPositionLots()
	if err != nil {
		return err
	}
	portfolio, err = bondportfolio.PopulateLotsWithBonds(portfolio)
	if err != nil {
		return err
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
				fmt.Fprint(sb, coupon.PerBondAmount)
				sb.WriteString(" на ")
				sb.WriteString(strconv.Itoa(int(portfolio[i].Quantity)))
				sb.WriteString(" шт. = ")
				sb.WriteString(forexservice.GetCurrencySymbol(portfolio[i].Bond.NominalCurrency))
				fmt.Fprint(sb, coupon.PerBondAmount*portfolio[i].Quantity)
			}
		}
	}

	return nil
}

func writeAnalysisToFile(sb *strings.Builder) error {
	fileName := "portfolio-analysis.txt"
	file, err := os.Create(fileName)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return err
	}
	defer file.Close()

	file.WriteString(sb.String())

	return nil
}
