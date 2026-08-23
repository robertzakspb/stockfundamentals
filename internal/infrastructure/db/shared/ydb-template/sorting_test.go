package ydbtemplate

import (
	"strings"
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_addSortingToQuery_Negative_NoColumn(t *testing.T) {
	columnName := ""
	direction := ">"
	var sb strings.Builder
	sb.WriteString("someQuery")

	addSortingToQuery(&sb, columnName, direction)

	//If no column was provided, it is expected that the query remains unaltered
	test.AssertEqual(t, "someQuery", sb.String())
}

func Test_addSortingToQuery_Negative_NoSortDirection(t *testing.T) {
	columnName := "age"
	direction := ""
	var sb strings.Builder
	sb.WriteString("SELECT * FROM user")

	expectedQuery := "SELECT * FROM user\nORDER BY age DESC\n"

	addSortingToQuery(&sb, columnName, direction)

	//If no column was provided, it is expected that the query remains unaltered
	test.AssertEqual(t, expectedQuery, sb.String())
}

func Test_addSortingToQuery_Positive_Ascending(t *testing.T) {
	columnName := "age"
	direction := ">"
	var sb strings.Builder
	sb.WriteString("SELECT * FROM user")

	expectedQuery := "SELECT * FROM user\nORDER BY age ASC\n"

	addSortingToQuery(&sb, columnName, direction)

	//If no column was provided, it is expected that the query remains unaltered
	test.AssertEqual(t, expectedQuery, sb.String())
}

func Test_addSortingToQuery_Positive_Descending(t *testing.T) {
	columnName := "age"
	direction := "<"
	var sb strings.Builder
	sb.WriteString("SELECT * FROM user")

	expectedQuery := "SELECT * FROM user\nORDER BY age DESC\n"

	addSortingToQuery(&sb, columnName, direction)

	//If no column was provided, it is expected that the query remains unaltered
	test.AssertEqual(t, expectedQuery, sb.String())
}
