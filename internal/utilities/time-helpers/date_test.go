package timehelpers

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_AreEqualDates_Positive(t *testing.T) {
	date1 := time.Now()
	date2 := time.Now().Add(time.Second)

	areEqualDates := AreEqualDates(date1, date2)
	test.AssertTrue(t, areEqualDates)
}

func Test_AreEqualDates_Negative(t *testing.T) {

	date1 := time.Now()
	date2 := time.Now().Add(time.Hour * 24)

	areEqualDates := AreEqualDates(date1, date2)
	test.AssertFalse(t, areEqualDates)
}

func Test_DateIsEarlierOrSameDate_SameDate(t *testing.T) {
	//Positive case where the same date is provided
	time1, time2 := time.Now(), time.Now()

	test.AssertTrue(t, DateIsEarlierOrSameDate(time1, time2))
}

func Test_DateIsEarlierOrSameDate_EarlierDate(t *testing.T) {
	//Positive case where an earlier date is provided
	time1, time2 := time.Now().AddDate(0, 0, -1), time.Now()

	test.AssertTrue(t, DateIsEarlierOrSameDate(time1, time2))
}

func Test_DateIsEarlierOrSameDate_LaterDate(t *testing.T) {
	//Negative case where a later date is provided
	time1, time2 := time.Now().AddDate(0, 0, 1), time.Now()

	test.AssertFalse(t, DateIsEarlierOrSameDate(time1, time2))
}

func Test_DateIsLaterOrSameDate_SameDate(t *testing.T) {
	//Positive case where the same date is provided
	time1, time2 := time.Now(), time.Now()

	test.AssertTrue(t, DateIsEarlierOrSameDate(time1, time2))
}

func Test_DateIsLaterOrSameDate_EaerlierDate(t *testing.T) {
	//Negative case where an earlier date is provided
	time1, time2 := time.Now().AddDate(0, 0, -1), time.Now()

	test.AssertFalse(t, DateIsLaterOrSameDate(time1, time2))
}

func Test_DateIsLaterOrSameDate_LaterDate(t *testing.T) {
	//Positive case where a later date is provided
	time1, time2 := time.Now().AddDate(0, 0, 1), time.Now()

	test.AssertTrue(t, DateIsLaterOrSameDate(time1, time2))
}

func Test_DateIsEarlier_SameDate(t *testing.T) {
	//Negative case where the same dates are provided
	time1, time2 := time.Now(), time.Now()

	test.AssertFalse(t, DateIsEarlier(time1, time2))
}

func Test_DateIsEarlier_LaterDate(t *testing.T) {
	//Negative case where a later date is provided
	time1, time2 := time.Now().AddDate(0, 0, 1), time.Now()

	test.AssertFalse(t, DateIsEarlier(time1, time2))
}

func Test_DateIsEarlier_EarlierDate(t *testing.T) {
	//Positive case where an earlier date is provided
	time1, time2 := time.Now().AddDate(0, 0, -1), time.Now()

	test.AssertTrue(t, DateIsEarlier(time1, time2))
}

func Test_DateIsLater_SameDate(t *testing.T) {
	//Negative case where the same dates are provided
	time1, time2 := time.Now(), time.Now()

	test.AssertFalse(t, DateIsLater(time1, time2))
}

func Test_DateIsLater_LaterDate(t *testing.T) {
	//Positive case where a later date is provided
	time1, time2 := time.Now().AddDate(0, 0, 1), time.Now()

	test.AssertTrue(t, DateIsLater(time1, time2))
}

func Test_DateIsLater_EarlierDate(t *testing.T) {
	//Negative case where an earlier date is provided
	time1, time2 := time.Now().AddDate(0, 0, -1), time.Now()

	test.AssertFalse(t, DateIsLater(time1, time2))
}
