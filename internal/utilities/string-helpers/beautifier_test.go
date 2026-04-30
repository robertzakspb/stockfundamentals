package stringhelpers

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_BeautifyNumber_LessThanThousand(t *testing.T) {
	number := 123.31

	beautified, err := BeautifyNumber(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "+123.31", beautified)
}

func Test_BeautifyNumber_LessThanThousandNegative(t *testing.T) {
	number := -123.31

	beautified, err := BeautifyNumber(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "-123.31", beautified)
}

func Test_BeautifyNumber_LessThanMillion(t *testing.T) {
	number := 123233.31

	beautified, err := BeautifyNumber(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "+123.23 тыс.", beautified)
}

func Test_BeautifyNumber_LessThanMillionNegative(t *testing.T) {
	number := -123000.31

	beautified, err := BeautifyNumber(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "-123 тыс.", beautified)
}
func Test_BeautifyNumber_LessThanBillion(t *testing.T) {
	number := 123233500.31

	beautified, err := BeautifyNumber(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "+123.23 млн.", beautified)
}

func Test_BeautifyNumber_LessThanBillionNegative(t *testing.T) {
	number := -123000405.31

	beautified, err := BeautifyNumber(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "-123 млн.", beautified)
}

func Test_BeautifyNumber_Zero(t *testing.T) {
	number := 0.0

	beautified, err := BeautifyNumber(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "+0", beautified)
}

func Test_BeatifyPercentage_Zero(t *testing.T) {
	percentage := 0.0

	beautified, err := BeatufityPercentage(percentage)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "+0.0%", beautified)
}

func Test_BeatifyPercentage_LessThan100(t *testing.T) {
	percentage := 0.5733

	beautified, err := BeatufityPercentage(percentage)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "+57.3%", beautified)
}

func Test_BeatifyPercentage_LessThan100Negative(t *testing.T) {
	percentage := -0.5733

	beautified, err := BeatufityPercentage(percentage)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "-57.3%", beautified)
}

func Test_BeatifyPercentage_MoreThan100(t *testing.T) {
	percentage := 2.5733

	beautified, err := BeatufityPercentage(percentage)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "x3.6", beautified)
}

func Test_BeatifyPercentage_MoreThan100Negative(t *testing.T) {
	percentage := -2.5733

	beautified, err := BeatufityPercentage(percentage)

	test.AssertNoError(t, err)
	test.AssertEqual(t, "-257.3%", beautified)
}
