package timehelpers

import "time"

func IsFutureDate(timestamp time.Time) bool {
	currentYear, currentMonth, currentDay := time.Now().Date()
	targetYear, targetMonth, targetDay := timestamp.Date()

	isFutureDate := targetYear >= currentYear &&
		targetMonth >= currentMonth &&
		targetDay > currentDay

	return isFutureDate
}

func DateIsToday(timestamp time.Time) bool {
	currentYear, currentMonth, currentDay := time.Now().Date()
	targetYear, targetMonth, targetDay := timestamp.Date()

	isToday := targetYear == currentYear &&
		targetMonth == currentMonth &&
		targetDay == currentDay

	return isToday
}

func IsPastDate(timestamp time.Time) bool {
	currentYear, currentMonth, currentDay := time.Now().Date()
	targetYear, targetMonth, targetDay := timestamp.Date()

	isPastDate := targetYear <= currentYear &&
		targetMonth <= currentMonth &&
		targetDay < currentDay

	return isPastDate
}

func IsTodayOrFutureDate(timestamp time.Time) bool {
	return IsFutureDate(timestamp) || DateIsToday(timestamp)
}

func IsTodayOrPastDate(timestamp time.Time) bool {
	return IsPastDate(timestamp) || DateIsToday(timestamp)
}
