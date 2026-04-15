package stringhelpers

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_RemoveDuplicatesFrom_EmptySlice(t *testing.T) {
	slice := []string{}

	noDuplicates := RemoveDuplicatesFrom(slice)

	test.AssertEqual(t, 0, len(noDuplicates))
}

func Test_RemoveDuplicatesFrom_NoDuplicates(t *testing.T) {
	slice := []string{"apple", "banana"}

	noDuplicates := RemoveDuplicatesFrom(slice)

	test.AssertEqual(t, 2, len(noDuplicates))
	test.AssertEqual(t, "apple", noDuplicates[0])
	test.AssertEqual(t, "banana", noDuplicates[1])
}

func Test_RemoveDuplicatesFrom_TwoDuplicates(t *testing.T) {
	slice := []string{"apple", "kiwi", "banana", "dates", "apple", "orange", "banana"}

	noDuplicates := RemoveDuplicatesFrom(slice)

	test.AssertEqual(t, 5, len(noDuplicates))
	test.AssertEqual(t, "apple", noDuplicates[0])
	test.AssertEqual(t, "kiwi", noDuplicates[1])
	test.AssertEqual(t, "banana", noDuplicates[2])
	test.AssertEqual(t, "dates", noDuplicates[3])
	test.AssertEqual(t, "orange", noDuplicates[4])
}
