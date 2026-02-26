package ydbfilter

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func Test_AddYqlVarDeclarations(t *testing.T) {
	filters := []YdbFilter{
		{"age", GreaterThan, types.Int64Value(20)},
		{"name", Equal, types.UTF8Value("Robert")},
		{"is_resident", Equal, types.BoolValue(false)},
	}

	expected := "DECLARE $age_filter AS Int64;\nDECLARE $name_filter AS Utf8;\nDECLARE $is_resident_filter AS Bool;\n"
	actual := AddYqlVarDeclarations(filters)

	test.AssertEqualStrings(t, expected, actual)
}

func Test_AddYqlVarDeclarations_NoFilters(t *testing.T) {
	filters := []YdbFilter{}

	expected := ""
	actual := AddYqlVarDeclarations(filters)

	test.AssertEqualStrings(t, expected, actual)
}

func Test_MakeWhereClause(t *testing.T) {
	filters := []YdbFilter{
		{"age", GreaterThan, types.Int64Value(20)},
		{"name", Equal, types.UTF8Value("Robert")},
		{"is_resident", Equal, types.BoolValue(false)},
	}

	expected := " WHERE\n age > $age_filter AND name = $name_filter AND is_resident = $is_resident_filter"
	actual := MakeWhereClause(filters)
	test.AssertEqualStrings(t, expected, actual)
}

func Test_MakeWhereClause_NoFilters(t *testing.T) {
	filters := []YdbFilter{}

	expected := ""
	actual := MakeWhereClause(filters)

	test.AssertEqualStrings(t, expected, actual)
}
