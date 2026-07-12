package ydbfilter

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_MapQueryFiltersToYdb_Negative_InvalidQuery(t *testing.T) {
	filters := map[string][]string{
		"name": {"=,"},
	}

	_, err := MapQueryFiltersToYdb[string](filters)

	test.AssertError(t, err)
}

func Test_MapQueryFiltersToYdb_Positive_SingleQueryStringEqual(t *testing.T) {
	type User struct {
		Name string `json:"name" sql:"db_name"`
	}
	filters := map[string][]string{
		"name": {"=,Robert"},
	}

	ydbFilters, err := MapQueryFiltersToYdb[User](filters)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 1, len(ydbFilters))
	test.AssertEqual(t, "db_name", ydbFilters[0].YqlColumnName)
	test.AssertEqual(t, Equal, ydbFilters[0].Condition)
	test.AssertEqual(t, "\"Robert\"u", ydbFilters[0].ConditionValue.Yql())
}
