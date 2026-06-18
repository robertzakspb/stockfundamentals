package ydbfilter

import (
	"fmt"

	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func Declare(yqlVarName string, yqlValue types.Value) string {
	if yqlVarName == "" {
		return ""
	}
	return fmt.Sprintf(
		"DECLARE %s AS %s;\n",
		yqlVarName,
		yqlValue.Type().Yql(),
	)
}

/*
Used to generate a sample filter name for yql queries. E.g. if a column is
titled "age", the appropriate filter variable name would be "$age_filter"
*/
func MakeColumnFilterName(columnName string, customPostfix string) string {
	if columnName == "" {
		return ""
	}
	return "$" + columnName + "_filter" + customPostfix
}

func groupFiltersByColumnName(filters []YdbFilter) map[string][]YdbFilter {
	groupedFilters := map[string][]YdbFilter{}

	for i := range filters {
		groupedFilters[filters[i].YqlColumnName] = append(groupedFilters[filters[i].YqlColumnName], filters[i])
	}

	return groupedFilters
}
