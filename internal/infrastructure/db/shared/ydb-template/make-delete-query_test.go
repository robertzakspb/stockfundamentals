package ydbtemplate

import (
	"strings"
	"testing"

	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
)

// Tests a query to delete all records in a particular table
func Test_generateDeleteQuery_DeleteAll(t *testing.T) {
	filters := []ydbfilter.YdbFilter{}
	tablePath := "`user/permission`"

	deleteQuery := generateDeleteQuery(filters, tablePath)

	test.AssertEqual(t, deleteQuery, "DELETE FROM `user/permission` ")
}

// Tests a query with a single filter
func Test_generateDeleteQuery_OneFilter(t *testing.T) {
	filters := []ydbfilter.YdbFilter{{
		YqlColumnName:  "is_enabled",
		Condition:      ydbfilter.Equal,
		ConditionValue: types.BoolValue(true),
	}}
	tablePath := "`user/permission`"

	deleteQuery := generateDeleteQuery(filters, tablePath)

	test.AssertEqual(t, deleteQuery, "DECLARE $is_enabled_filter1 AS Bool;\nDELETE FROM `user/permission`  WHERE\n is_enabled = $is_enabled_filter1")
}

// Tests a query with two filters
func Test_generateDeleteQuery_TwoFilters(t *testing.T) {
	filters := []ydbfilter.YdbFilter{
		{
			YqlColumnName:  "is_enabled",
			Condition:      ydbfilter.Equal,
			ConditionValue: types.BoolValue(true),
		},
		{
			YqlColumnName:  "age",
			Condition:      ydbfilter.GreaterThanOrEqualTo,
			ConditionValue: types.Int64Value(18),
		}}
	tablePath := "`user/permission`"

	deleteQuery := generateDeleteQuery(filters, tablePath)

	baseQuery := "DECLARE $is_enabled_filter1 AS Bool;\nDECLARE $age_filter1 AS Int64;\nDELETE FROM `user/permission`  WHERE\n "
	firstWhereParameter := "is_enabled = $is_enabled_filter1"
	secondWhereParameter := "age >= $age_filter1"

	test.AssertTrue(t, strings.Contains(deleteQuery, baseQuery))
	test.AssertTrue(t, strings.Contains(deleteQuery, firstWhereParameter))
	test.AssertTrue(t, strings.Contains(deleteQuery, secondWhereParameter))
}
