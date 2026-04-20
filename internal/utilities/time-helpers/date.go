package timehelpers

import "time"

func AreEqualDates(time1, time2 time.Time) bool {
	time1Year, time1Month, time1Day := time1.Date()
	time2Year, time2Month, time2Day := time2.Date()

	areEqual := time1Year == time2Year && time1Month == time2Month && time1Day == time2Day

	return areEqual
}

// Indicates whether time1 is the same date or an earlier date than time2
func DateIsEarlierOrSameDate(time1, time2 time.Time) bool {
	if AreEqualDates(time1, time2) {
		return true
	}

	date1 := ConvertTimeToMidnightUTC(time1)
	date2 := ConvertTimeToMidnightUTC(time2)

	return date1.Before(date2)
}

// Indicates whether time1 is an earlier date than time2.
// If the same dates are provided, the function returns false
func DateIsEarlier(time1, time2 time.Time) bool {
	return !DateIsLaterOrSameDate(time1, time2)
}

// Indicates whether time1 is the same date or a later date than time2
func DateIsLaterOrSameDate(time1, time2 time.Time) bool {
	if AreEqualDates(time1, time2) {
		return true
	}

	date1 := ConvertTimeToMidnightUTC(time1)
	date2 := ConvertTimeToMidnightUTC(time2)

	return date1.After(date2)
}

// Indicates whether time1 is a later date than time2.
// If the same dates are provided, the function returns false
func DateIsLater(time1, time2 time.Time) bool {
	return !DateIsEarlierOrSameDate(time1, time2)
}

// Converts a timestamp like "Jan 2, 5:23 PM" to "Jan 2, 00:00" in the UTC time zone
func ConvertTimeToMidnightUTC(timestamp time.Time) time.Time {
	year, month, day := timestamp.Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	return date
}
