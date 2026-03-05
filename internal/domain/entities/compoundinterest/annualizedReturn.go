package compoundinterest

import (
	"math"
	"time"
)

// Note: 12% is returned as 0.12
func CalcAnnualizedReturn(cumulativeReturn float64, startDate, endDate time.Time) float64 {
	if cumulativeReturn == 0 || startDate.Equal(endDate) || startDate.After(endDate) {
		return 0
	}

	daysHeld := endDate.Sub(startDate).Hours() / 24
	annualizedReturn := math.Pow(1+cumulativeReturn, 365/daysHeld) - 1

	return annualizedReturn
}
