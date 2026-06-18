package ydbfilter

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func Test_Declare_EmptyString(t *testing.T) {
	expected := ""
	actual := Declare("", types.UTF8Value(""))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_Timestamp(t *testing.T) {
	expected := "DECLARE payout_date AS Datetime;\n"
	actual := Declare("payout_date", types.DatetimeValueFromTime(time.Now()))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_String(t *testing.T) {
	expected := "DECLARE name AS Utf8;\n"
	actual := Declare("name", types.UTF8Value("Robert"))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_Int64(t *testing.T) {
	expected := "DECLARE age AS Int64;\n"
	actual := Declare("age", types.Int64Value(10))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_Int32(t *testing.T) {
	expected := "DECLARE age AS Int32;\n"
	actual := Declare("age", types.Int32Value(10))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_Int16(t *testing.T) {
	expected := "DECLARE age AS Int16;\n"
	actual := Declare("age", types.Int16Value(10))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_UInt8(t *testing.T) {
	expected := "DECLARE age AS Uint8;\n"
	actual := Declare("age", types.Uint8Value(10))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_Uint64(t *testing.T) {
	expected := "DECLARE age AS Uint64;\n"
	actual := Declare("age", types.Uint64Value(10))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_Uint32(t *testing.T) {
	expected := "DECLARE age AS Uint32;\n"
	actual := Declare("age", types.Uint32Value(10))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_Uint16(t *testing.T) {
	expected := "DECLARE age AS Uint16;\n"
	actual := Declare("age", types.Uint16Value(10))

	test.AssertEqual(t, expected, actual)
}

func Test_Declare_Uint8(t *testing.T) {
	expected := "DECLARE age AS Uint8;\n"
	actual := Declare("age", types.Uint8Value(10))

	test.AssertEqual(t, expected, actual)
}

func Test_MakeColumnFilterNameWithEmptyString(t *testing.T) {
	expected := ""
	actual := MakeColumnFilterName("", "")

	test.AssertEqual(t, expected, actual)
}

func Test_MakeColumnFilterName(t *testing.T) {
	expected := "$age_filter"
	actual := MakeColumnFilterName("age", "")

	test.AssertEqual(t, expected, actual)
}

func Test_groupFiltersByColumnName_Positive(t *testing.T) {
	filters := []YdbFilter{
		{
			YqlColumnName: "name1",
		},
		{
			YqlColumnName: "name2",
		},
		{
			YqlColumnName: "name2",
		},
		{
			YqlColumnName: "name3",
		},
		{
			YqlColumnName: "name1",
		},
		{
			YqlColumnName: "name1",
		},
	}
	grouped := groupFiltersByColumnName(filters)

	test.AssertEqual(t, 3, len(grouped))
	test.AssertEqual(t, 3, len(grouped["name1"]))
	test.AssertEqual(t, 2, len(grouped["name2"]))
	test.AssertEqual(t, 1, len(grouped["name3"]))
}

func Test_groupFiltersByColumnName_EmptySlice(t *testing.T) {
	filters := []YdbFilter{}
	grouped := groupFiltersByColumnName(filters)

	test.AssertEqual(t, 0, len(grouped))
}
