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

	expected := "DECLARE $age_filter1 AS Int64;\nDECLARE $name_filter1 AS Utf8;\nDECLARE $is_resident_filter1 AS Bool;\n"
	actual := AddYqlVarDeclarations(filters)

	test.AssertEqual(t, expected, actual)
}

func Test_AddYqlVarDeclarations_NoFilters(t *testing.T) {
	filters := []YdbFilter{}

	expected := ""
	actual := AddYqlVarDeclarations(filters)

	test.AssertEqual(t, expected, actual)
}

func Test_MakeWhereClause(t *testing.T) {
	filters := []YdbFilter{
		{"age", GreaterThan, types.Int64Value(20)},
		{"name", Equal, types.UTF8Value("Robert")},
		{"is_resident", Equal, types.BoolValue(false)},
	}

	expected := " WHERE\n age > $age_filter1 AND name = $name_filter1 AND is_resident = $is_resident_filter1"
	actual := MakeWhereClause(filters)
	test.AssertEqual(t, expected, actual)
}

func Test_MakeWhereClause_NoFilters(t *testing.T) {
	filters := []YdbFilter{}

	expected := ""
	actual := MakeWhereClause(filters)

	test.AssertEqual(t, expected, actual)
}

func Test_convertQueryParamToYdbFilter_ListOFStrings(t *testing.T) {
	jsonParameter := "accountId"
	type Foo struct {
		AccountId string `json:"accountId" sql:"account_id"`
	}
	queryValues := []string{"IN", "id1", "id2"}

	filter, err := convertQueryParamToYdbFilter(jsonParameter, Foo{}, queryValues)

	test.AssertNoError(t, err)
	test.AssertEqual(t, Contains, filter.Condition)
	test.AssertEqual(t, "account_id", filter.YqlColumnName)
	test.AssertEqual(t, "List<Utf8>", filter.ConditionValue.Type().String())
}

func Test_convertQueryParamToYdbFilter_SingleInteger(t *testing.T) {
	jsonParameter := "quantity"
	type Foo struct {
		Quantity string `json:"quantity" sql:"quantity_db"`
	}
	queryValues := []string{"=", "5"}

	filter, err := convertQueryParamToYdbFilter(jsonParameter, Foo{}, queryValues)

	test.AssertNoError(t, err)
	test.AssertEqual(t, Equal, filter.Condition)
	test.AssertEqual(t, "quantity_db", filter.YqlColumnName)
	test.AssertEqual(t, "Double", filter.ConditionValue.Type().String())
}
