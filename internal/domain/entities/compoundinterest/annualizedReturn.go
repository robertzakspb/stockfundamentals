package compoundinterest

import (
	"math"
	"time"
)

// Note: 12% is returned as 0.12
func CalcAnnualizedReturn(totalReturnPercentage float64, startDate, endDate time.Time) float64 {
	if totalReturnPercentage == 0 || startDate.Equal(endDate) || startDate.After(endDate) {
		return -1
	}

	daysHeld := endDate.Sub(startDate).Hours() / 24
	annualizedReturn := math.Pow(1+totalReturnPercentage, 365/daysHeld) - 1

	return annualizedReturn
}
