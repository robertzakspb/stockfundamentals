package timehelpers

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_ParseMonth_Negative_NonExistentMonth(t *testing.T) {
	nonExistentMonth := "Blabruary"

	_, err := ParseMonth(nonExistentMonth)

	test.AssertError(t, err)
}

func Test_ParseMonth_Positive_January(t *testing.T) {
	monthString := "January"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 1, month)
}

func Test_ParseMonth_Positive_February(t *testing.T) {
	monthString := "February"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 2, month)
}

func Test_ParseMonth_Positive_March(t *testing.T) {
	monthString := "March"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 3, month)
}

func Test_ParseMonth_Positive_April(t *testing.T) {
	monthString := "April"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 4, month)
}

func Test_ParseMonth_Positive_May(t *testing.T) {
	monthString := "May"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 5, month)
}

func Test_ParseMonth_Positive_June(t *testing.T) {
	monthString := "June"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 6, month)
}

func Test_ParseMonth_Positive_July(t *testing.T) {
	monthString := "July"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 7, month)
}

func Test_ParseMonth_Positive_August(t *testing.T) {
	monthString := "August"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 8, month)
}

func Test_ParseMonth_Positive_September(t *testing.T) {
	monthString := "September"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 9, month)
}

func Test_ParseMonth_Positive_October(t *testing.T) {
	monthString := "October"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 10, month)
}

func Test_ParseMonth_Positive_November(t *testing.T) {
	monthString := "November"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 11, month)
}

func Test_ParseMonth_Positive_December(t *testing.T) {
	monthString := "December"

	month, err := ParseMonth(monthString)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 12, month)
}



