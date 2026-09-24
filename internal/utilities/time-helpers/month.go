package timehelpers

import (
	"errors"
	"strings"
	"time"
)

func ParseMonth(month string) (time.Month, error) {
	switch strings.ToUpper(month) {
	case "JANUARY":
		return 1, nil
	case "FEBRUARY":
		return 2, nil
	case "MARCH":
		return 3, nil
	case "APRIL":
		return 4, nil
	case "MAY":
		return 5, nil
	case "JUNE":
		return 6, nil
	case "JULY":
		return 7, nil
	case "AUGUST":
		return 8, nil
	case "SEPTEMBER":
		return 9, nil
	case "OCTOBER":
		return 10, nil
	case "NOVEMBER":
		return 11, nil
	case "DECEMBER":
		return 12, nil
	default:
		return time.January, errors.New("Failed to parse the month from the provided string: " + month)
	}
}
