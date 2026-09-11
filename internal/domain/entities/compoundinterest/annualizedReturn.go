package compoundinterest

import (
	"errors"
	"math"
	"time"
)

// Note: 12% is returned as 0.12
func CalcAnnualizedReturn(totalReturnPercentage float64, startDate, endDate time.Time) (float64, error) {
	if totalReturnPercentage == 0 || startDate.Equal(endDate) || startDate.After(endDate) {
		return -1, errors.New("Invalid data was provided to the CalcAnnualizedReturn functon")
	}

	daysHeld := endDate.Sub(startDate).Hours() / 24
	if int(daysHeld) == 0 {
		return -1, errors.New("Unable to calculate the return, as the position has been held for 0 days")
	}
	annualizedReturn := math.Pow(1+totalReturnPercentage, 365/daysHeld) - 1

	if math.IsInf(annualizedReturn, 0) {
		return -1, errors.New("Found infinity when calculating the annualized return")
	}

	return annualizedReturn, nil
}
