package ydbfilter

import (
	"fmt"

	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func Declare(yqlVarName string, yqlValue types.Value) string {
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
func MakeColumnFilterName(columnName string) string {
	return "$" + columnName + "_filter"
}
