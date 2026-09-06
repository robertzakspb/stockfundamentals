package compoundinterest

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

const initialGuess = 0.2
const accuracy = 0.00000001
const maxIterations = 50

//Nice write-up on IRR: https://medium.com/@_orcaman/package-financial-for-golang-the-math-behind-the-irr-function-1eedf225d9f
//XIRR formula in LibreOffice: https://wiki.openoffice.org/wiki/Documentation/How_Tos/Calc:_XIRR_function

func InternalRateOfReturn(cashflows []float64, dates []time.Time) (float64, error) {
	irr, err := internalRateOfReturn(cashflows, dates)

	return irr, err
}

// Calculates the internal rate of return for the proided cash flows and their corresponding dates
func internalRateOfReturn(cashflows []float64, dates []time.Time) (float64, error) {
	if len(cashflows) == 0 {
		return -1, errors.New("IRR cannot be calculated due to the missing cash flows")
	}
	if len(dates) != len(cashflows) {
		return -1, errors.New("Date and cashflow counts must be identical")
	}

	x0 := initialGuess
	var x1 float64

	for i := range maxIterations {
		var fValue, fDerivative float64
		for k := range cashflows {
			fValue += fx_xirr(cashflows[k], x0, dates[k], dates[0])
			fDerivative += fx_d_xirr(cashflows[k], x0, dates[k], dates[0])
		}
		if fDerivative == 0 {
			continue
		}
		x1 = x0 - fValue/fDerivative

		if math.Abs(x1-x0) < accuracy {
			fmt.Println("Iteration count: ", strconv.Itoa(i))
			return x1, nil
		}
		x0 = x1
	}

	return -1, errors.New("Failed to find the internal rate of return")
}

// Calculates the value of the XIRR function given the cashflow, its date, and the initial date
func fx_xirr(cashflow, x0 float64, cashflowDate, startDate time.Time) float64 {
	daysSinceStartUntilCashflow := startDate.Sub(cashflowDate).Hours() / 24
	fx := cashflow * math.Pow(1.0+x0, daysSinceStartUntilCashflow/365.0)
	return fx
}

// Calculates the deriviative of the function
func fx_d_xirr(cashflow, x0 float64, cashflowDate, startDate time.Time) float64 {
	daysSinceStartUntilCashflow := startDate.Sub(cashflowDate).Hours() / 24
	dfx := (1.0 / 365.0) * (daysSinceStartUntilCashflow) * cashflow * math.Pow((x0+1.0), ((daysSinceStartUntilCashflow/365.0)-1.0))
	return dfx
}

// Converts an IRR calculated for non-annual cashflows (e.g. monthly, quarterly) into an annualized IRR
func annualizeIrrRate(irr float64, cashflowCountPerYear int) float64 {
	//One annual casfhlow means it's already an annualized IRR and we can hence return it
	if cashflowCountPerYear == 1 {
		return irr
	}

	annualIRR := math.Pow((1+irr), float64(cashflowCountPerYear)) - 1
	return annualIRR
}
