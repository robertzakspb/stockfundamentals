package typeconverter

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_GetFloat_Float64(t *testing.T) {
	const number float64 = 45.7

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, number)
}

func Test_GetFloat_Float32(t *testing.T) {
	const number float32 = 45.7

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Int64(t *testing.T) {
	const number int64 = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Int32(t *testing.T) {
	const number int32 = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Int16(t *testing.T) {
	const number int16 = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Int8(t *testing.T) {
	const number int8 = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Int(t *testing.T) {
	const number int = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Uint64(t *testing.T) {
	const number uint64 = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Uint32(t *testing.T) {
	const number uint32 = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Uint16(t *testing.T) {
	const number uint16 = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Uint8(t *testing.T) {
	const number uint8 = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_Uint(t *testing.T) {
	const number uint = 45

	parsedFloat, err := GetFloat(number)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, float64(number))
}

func Test_GetFloat_ValidStringPositiveNumber(t *testing.T) {
	const str = "45.7"

	parsedFloat, err := GetFloat(str)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, 45.7)
}

func Test_GetFloat_ValidStringNegativeNumber(t *testing.T) {
	const str = "-45.7"

	parsedFloat, err := GetFloat(str)

	test.AssertNoError(t, err)
	test.AssertEqual(t, parsedFloat, -45.7)
}

func Test_GetFloat_InvalidString(t *testing.T) {
	const str = "abc"

	_, err := GetFloat(str)

	test.AssertError(t, err)
}
