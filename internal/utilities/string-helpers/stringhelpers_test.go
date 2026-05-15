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

func Test_SplitInBatchesOf_NoElements(t *testing.T) {
	emptySlice := []string{}

	batches := SplitInBatchesOf(4, emptySlice)

	test.AssertEqual(t, 1, len(batches))
	test.AssertEqual(t, 0, len(batches[0]))
}

func Test_SplitInBatchesOf_SizeExceedsSliceLength(t *testing.T) {
	slice := []string{"test1", "test2"}

	batches := SplitInBatchesOf(4, slice)

	test.AssertEqual(t, "test1", batches[0][0])
	test.AssertEqual(t, "test2", batches[0][1])
	test.AssertEqual(t, 1, len(batches))
}

func Test_SplitInBatchesOf_OneBatch(t *testing.T) {
	emptySlice := []string{"test1", "test2"}

	batches := SplitInBatchesOf(2, emptySlice)

	test.AssertEqual(t, 1, len(batches))
	test.AssertEqual(t, "test1", batches[0][0])
	test.AssertEqual(t, "test2", batches[0][1])
}

func Test_SplitInBatchesOf_MultipleEquallySizedBatches(t *testing.T) {
	slice := []string{"test1", "test2", "test3", "test4", "test5", "test6"}

	batches := SplitInBatchesOf(3, slice)

	test.AssertEqual(t, 2, len(batches))
	test.AssertEqual(t, "test1", batches[0][0])
	test.AssertEqual(t, "test2", batches[0][1])
	test.AssertEqual(t, "test3", batches[0][2])
	test.AssertEqual(t, "test4", batches[1][0])
	test.AssertEqual(t, "test5", batches[1][1])
	test.AssertEqual(t, "test6", batches[1][2])
}

func Test_SplitInBatchesOf_MultipleBatchesWithLastUnfilledBatch(t *testing.T) {
	slice := []string{"test1", "test2", "test3", "test4", "test5", "test6", "test7", "test8"}

	batches := SplitInBatchesOf(3, slice)

	test.AssertEqual(t, 3, len(batches))
	test.AssertEqual(t, "test1", batches[0][0])
	test.AssertEqual(t, "test2", batches[0][1])
	test.AssertEqual(t, "test3", batches[0][2])
	test.AssertEqual(t, "test4", batches[1][0])
	test.AssertEqual(t, "test5", batches[1][1])
	test.AssertEqual(t, "test6", batches[1][2])
	test.AssertEqual(t, "test7", batches[2][0])
	test.AssertEqual(t, "test8", batches[2][1])
}
