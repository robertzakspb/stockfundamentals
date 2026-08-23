package ydbtemplate

import (
	"strings"
	"testing"

	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

func Test_GenerateGetQuery_NoFilters(t *testing.T) {
	type Foo struct {
		a string  `sql:"a"`
		b float64 `sql:"b"`
		c bool
	}
	tableName := "fooStore"
	expectedQuery := "SELECT a, b, FROM fooStore "

	query, err := generateGetQuery[Foo]([]ydbfilter.YdbFilter{}, "", "", tableName, 0)

	test.AssertNoError(t, err)
	test.AssertEqual(t, expectedQuery, query)
}

func Test_GenerateGetQuery_MissingSqlTag(t *testing.T) {
	type Foo struct {
		a string  `sql:"a"`
		b float64 `sql:"b"`
		c bool
	}

	tableName := "fooStore"
	filters := []ydbfilter.YdbFilter{
		{
			YqlColumnName:  "nonExistentTag",
			Condition:      ydbfilter.GreaterThan,
			ConditionValue: types.TextValue("test"),
		},
	}

	_, err := generateGetQuery[Foo](filters, "", "", tableName, 0)

	test.AssertError(t, err)
}

func Test_GenerateGetQuery_TwoFilters(t *testing.T) {
	type Foo struct {
		a string  `sql:"a"`
		b float64 `sql:"b"`
		c bool
	}
	firstDeclaration := "DECLARE $a_filter1 AS Utf8;\n"
	secondDeclaration := "DECLARE $b_filter1 AS Double;\n"
	selectPart := "SELECT a, b, FROM fooStore  WHERE\n"
	firstWhereArgument := "a > $a_filter1"
	secondWhereArgument := "b <= $b_filter1"


	tableName := "fooStore"
	filters := []ydbfilter.YdbFilter{
		{
			YqlColumnName:  "a",
			Condition:      ydbfilter.GreaterThan,
			ConditionValue: types.TextValue("test"),
		},
		{
			YqlColumnName:  "b",
			Condition:      ydbfilter.LessThanOrEqualTo,
			ConditionValue: types.DoubleValue(6.7),
		},
	}

	query, err := generateGetQuery[Foo](filters, "", "", tableName, 0)

	test.AssertNoError(t, err)
	test.AssertTrue(t, strings.Contains(query, firstDeclaration))
	test.AssertTrue(t, strings.Contains(query, secondDeclaration))
	test.AssertTrue(t, strings.Contains(query, selectPart))
	test.AssertTrue(t, strings.Contains(query, firstWhereArgument))
	test.AssertTrue(t, strings.Contains(query, secondWhereArgument))
}

func Test_generatePaginatedGetQuery_Negative_NoSortByParameter(t *testing.T) {
	type Foo struct{}
	sortBy := ""
	pageSize := 5

	_, err := generateGetQuery[Foo]([]ydbfilter.YdbFilter{}, sortBy, "", "", pageSize)

	test.AssertError(t, err)
}

func Test_generatePaginatedGetQuery_Positive_OnlySorting(t *testing.T) {
	type Foo struct {
		Age  int    `sql:"age"`
		Name string `sql:"name"`
	}
	sortBy := "age"
	sortDirection := "<"
	pageSize := 0
	tableName := "user"
	expectedQuery := "SELECT age, name, FROM user \nORDER BY age DESC\n"

	query, err := generateGetQuery[Foo]([]ydbfilter.YdbFilter{}, sortBy, sortDirection, tableName, pageSize)

	test.AssertNoError(t, err)
	test.AssertEqual(t, expectedQuery, query)
}

func Test_generatePaginatedGetQuery_Positive_SortingAndPagination(t *testing.T) {
	type Foo struct {
		Age  int    `json:"age" sql:"age"`
		Name string `sql:"name"`
	}
	sortBy := "age"
	sortDirection := "<"
	pageSize := 5
	tableName := "user"
	expectedQuery := "$pageSize = 5;\nSELECT age, name, FROM user \nORDER BY age DESC\n\nLIMIT $pageSize;"

	query, err := generateGetQuery[Foo]([]ydbfilter.YdbFilter{}, sortBy, sortDirection, tableName, pageSize)

	test.AssertNoError(t, err)
	test.AssertEqual(t, expectedQuery, query)
}
