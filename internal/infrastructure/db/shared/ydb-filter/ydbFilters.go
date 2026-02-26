package ydbfilter

import (
	"strings"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type YdbFilter struct {
	YqlColumnName  string
	Condition      YdbFilterCondition
	ConditionValue types.Value
}

type YdbFilterCondition string

const (
	GreaterThan          YdbFilterCondition = ">"
	GreaterThanOrEqualTo YdbFilterCondition = ">="
	LessThan             YdbFilterCondition = "<"
	LessThanOrEqualTo    YdbFilterCondition = "<="
	Contains             YdbFilterCondition = "IN"
	Equal                YdbFilterCondition = "="
)

// This function only sets the values of query parameters themeselves but does not add them WHERE (done by MakeWhereClause())
func SetQueryParams(filters []YdbFilter) *table.QueryParameters {
	params := []table.ParameterOption{}

	for _, filter := range filters {
		param := table.ValueParam(MakeColumnFilterName(filter.YqlColumnName), filter.ConditionValue)
		params = append(params, param)
	}

	return table.NewQueryParameters(params...)
}

func MakeWhereClause(filters []YdbFilter) string {
	if len(filters) == 0 {
		return ""
	}

	b := strings.Builder{}

	b.WriteString(" WHERE \n")

	for i, filter := range filters {
		b.WriteString(filter.YqlColumnName)
		b.WriteString(" ")
		b.WriteString(string(filter.Condition))
		b.WriteString(" ")
		b.WriteString(MakeColumnFilterName(filter.YqlColumnName))
		b.WriteString(" ")
		if i < len(filters)-1 {
			b.WriteString("AND")
			b.WriteString(" ")
		}
	}

	return b.String()
}

func AddYqlVarDeclarations(filters []YdbFilter) string {
	if len(filters) == 0 {
		return ""
	}

	b := strings.Builder{}

	for _, filter := range filters {
		yqlVarName := MakeColumnFilterName(filter.YqlColumnName)
		b.WriteString(Declare(yqlVarName, filter.ConditionValue))
	}

	return b.String()
}
