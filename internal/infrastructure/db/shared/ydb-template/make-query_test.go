package ydbtemplate

import (
	"testing"

	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func Test_GenerateGetQuery_NoFilters(t *testing.T) {
	type Foo struct {
		a string  `sql:"a"`
		b float64 `sql:"b"`
		c bool
	}
	tableName := "fooStore"
	expectedQuery := "SELECT a, b, FROM fooStore "

	query, err := generateGetQuery[Foo]([]ydbfilter.YdbFilter{}, tableName)

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

	_, err := generateGetQuery[Foo](filters, tableName)

	test.AssertError(t, err)
}

func Test_GenerateGetQuery_TwoFilters(t *testing.T) {
	type Foo struct {
		a string  `sql:"a"`
		b float64 `sql:"b"`
		c bool
	}
	expectedQuery := "DECLARE $a_filter1 AS Utf8;\nDECLARE $b_filter1 AS Double;\nSELECT a, b, FROM fooStore  WHERE\n a > $a_filter1 AND b <= $b_filter1"
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

	query, err := generateGetQuery[Foo](filters, tableName)

	test.AssertNoError(t, err)
	test.AssertEqual(t, expectedQuery, query)
}
