package timehelpers

// import (
// 	"testing"
// 	"time"

// 	"github.com/compoundinvest/stockfundamentals/internal/test"
// )

// func Test_DateIsToday(t *testing.T) {
// 	today := time.Now()

// 	test.AssertEqual(t, true, DateIsToday(today))
// }

// func Test_DateIsToday_Negative(t *testing.T) {
// 	threeDaysAgo := time.Now().Add(-time.Hour * 72)

// 	test.AssertEqual(t, false, DateIsToday(threeDaysAgo))
// }

// func Test_DateIsInFuture(t *testing.T) {
// 	threeDaysFromNow := time.Now().Add(time.Hour * 72)

// 	test.AssertEqual(t, true, IsFutureDate(threeDaysFromNow))
// }

// func Test_DateIsInFuture_Negative(t *testing.T) {
// 	threeDaysFromNow := time.Now().Add(time.Hour * 72)

// 	test.AssertEqual(t, true, IsFutureDate(threeDaysFromNow))
// }

// func Test_DateIsInPast(t *testing.T) {
// 	threeDaysAgo := time.Now().Add(-time.Hour * 72)

// 	test.AssertEqual(t, true, IsPastDate(threeDaysAgo))
// }

// func Test_DateIsInPast_Negative(t *testing.T) {
// 	threeDaysAgo := time.Now().Add(-time.Hour * 72)

// 	test.AssertEqual(t, true, IsPastDate(threeDaysAgo))
// }

// func Test_DateIsTodayOrInFuture(t *testing.T) {
// 	threeDaysFromNow := time.Now().Add(time.Hour * 72)

// 	test.AssertEqual(t, true, IsTodayOrFutureDate(threeDaysFromNow))
// }

// func Test_DateIsTodayOrInFuture_Negative(t *testing.T) {
// 	threeDaysAgo := time.Now().Add(-time.Hour * 72)

// 	test.AssertEqual(t, false, IsTodayOrFutureDate(threeDaysAgo))
// }

// func Test_DateIsTodayOrInPast(t *testing.T) {
// 	threeDaysAgo := time.Now().Add(-time.Hour * 72)

// 	test.AssertEqual(t, true, IsTodayOrPastDate(threeDaysAgo))
// }

// func Test_DateIsTodayOrInPast_Negative(t *testing.T) {
// 	threeDaysFromNow := time.Now().Add(time.Hour * 72)

// 	test.AssertEqual(t, false, IsTodayOrPastDate(threeDaysFromNow))
// }
