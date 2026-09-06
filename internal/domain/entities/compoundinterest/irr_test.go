package compoundinterest

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_InternalRateOfReturn_GovernmentBond(t *testing.T) {
	nominalPrice := 1000.0
	marketPrice := 848.7
	coupon := 33.41
	maturityDate := time.Date(2029, 3, 14, 0, 0, 0, 0, time.UTC)
	acruedInterest := 31.94
	cashFlows := []float64{-(marketPrice + acruedInterest), coupon, coupon, coupon, coupon, coupon, coupon, nominalPrice}

	dates := []time.Time{}
	dates = append(dates, time.Now())
	dates = append(dates, time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC))
	dates = append(dates, time.Date(2027, 3, 17, 0, 0, 0, 0, time.UTC))
	dates = append(dates, time.Date(2027, 9, 15, 0, 0, 0, 0, time.UTC))
	dates = append(dates, time.Date(2028, 3, 15, 0, 0, 0, 0, time.UTC))
	dates = append(dates, time.Date(2028, 9, 13, 0, 0, 0, 0, time.UTC))
	dates = append(dates, maturityDate)
	dates = append(dates, maturityDate)

	irr, err := InternalRateOfReturn(cashFlows, dates)

	test.AssertNoError(t, err)
	test.AssertEqualFloat(t, 0.1452, irr, 0.0001)
}

func Test_annualizeIrrRate_Positive_SemiAnnual(t *testing.T) {
	semiAnnualIrr := 0.06

	annualized := annualizeIrrRate(semiAnnualIrr, 2)

	test.AssertEqualFloat(t, 0.1236, annualized, 0.00001)
}

func Test_annualizeIrrRate_Positive_Annual(t *testing.T) {
	annualIrr := 0.06

	annualized := annualizeIrrRate(annualIrr, 1)

	test.AssertEqualFloat(t, 0.06, annualized, 0.00001)
}
