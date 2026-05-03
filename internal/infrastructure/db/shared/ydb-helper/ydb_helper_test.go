package ydbhelper

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
	"github.com/google/uuid"
)

func Test_ConvertUUIDsToYdbList(t *testing.T) {
	uuids := make([]uuid.UUID, 2)
	uuids[0] = uuid.MustParse("01d7a5d1-014a-41dc-80b6-af2fe480ae08")
	uuids[1] = uuid.MustParse("b58c7bea-0c00-4e8b-baa9-157abf07e1b8")

	ydbUUIDs := ConvertUUIDsToYdbList(uuids)

	test.AssertEqual(t, "[Uuid(\"01d7a5d1-014a-41dc-80b6-af2fe480ae08\"),Uuid(\"b58c7bea-0c00-4e8b-baa9-157abf07e1b8\")]", ydbUUIDs.Yql())
}

func Test_ConvertStringsToYdbList(t *testing.T) {
	strings := []string{"apple", "banana", "kiwi"}

	ydbStrings := ConvertStringsToYdbList(strings)

	test.AssertEqual(t, "[\"apple\"u,\"banana\"u,\"kiwi\"u]", ydbStrings.Yql())
}

func Test_ConverTimestampsToYdbDates(t *testing.T) {
	date1, _ := timehelpers.DateFromISOstring("2026-05-03")
	date2, _ := timehelpers.DateFromISOstring("2026-05-10")

	timestamps := []time.Time{date1, date2}

	ydbDates := ConvertTimestampsToYdbDates(timestamps)

	test.AssertEqual(t, "[Date(\"2026-05-03\"),Date(\"2026-05-10\")]", ydbDates.Yql())
}

func Test_ConvertToYdbDateTime(t *testing.T) {
	today, _ := timehelpers.DateFromISOstring("2026-05-03")

	ydbDate := ConvertToYdbDateTime(today)

	test.AssertEqual(t, "Datetime(\"2026-05-03T00:00:00Z\")", ydbDate.Yql())
}

func Test_ConvertToYdbDate(t *testing.T) {
	today, _ := timehelpers.DateFromISOstring("2026-05-03")

	ydbDate := ConvertToYdbDate(today)

	test.AssertEqual(t, "Date(\"2026-05-03\")", ydbDate.Yql())
}

func Test_ConvertToOptionalYDBdate(t *testing.T) {
	zeroTime := time.Time{}

	ydbDate := ConvertToOptionalYDBdate(zeroTime)

	test.AssertEqual(t, "Nothing(Optional<Date>)", ydbDate.Yql())
}
