package ydbtemplate

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_addSortingToQuery_Negative_NoColumn(t *testing.T) {
	columnName := ""
	yql := "someQuery"

	queryWithSorting := addSortingToQuery(yql, columnName)

	//If no column was provided, it is expected that the query remains unaltered
	test.AssertEqual(t, yql, queryWithSorting)
}


func Test_addSortingToQuery_Positive_Standard(t *testing.T) {
	columnName := "age"
	yql := "SELECT * FROM user"
	expectedQuery := "SELECT * FROM user\nORDER BY age\n"

	queryWithSorting := addSortingToQuery(yql, columnName)

	//If no column was provided, it is expected that the query remains unaltered
	test.AssertEqual(t, expectedQuery, queryWithSorting)
}
