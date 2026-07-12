package ydbfilter

import (
	"strconv"
	"strings"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/types"
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
	NotEqual             YdbFilterCondition = "!="
	Like                 YdbFilterCondition = "LIKE"
)

var ydbConditions = map[string]YdbFilterCondition{
	">":    GreaterThan,
	">=":   GreaterThanOrEqualTo,
	"<":    LessThan,
	"<=":   LessThanOrEqualTo,
	"IN":   Contains,
	"=":    Equal,
	"!=":   NotEqual,
	"LIKE": Like,
}

// This function only sets the values of query parameters themeselves but does not add them WHERE (done by AddWhereClause())
func SetQueryParams(filters []YdbFilter) *table.QueryParameters {
	params := []table.ParameterOption{}

	groupedFilters := groupFiltersByColumnName(filters)
	for _, filterList := range groupedFilters {
		for i, filter := range filterList {
			param := table.ValueParam(MakeColumnFilterName(filter.YqlColumnName, strconv.Itoa(i+1)), filter.ConditionValue)
			params = append(params, param)
		}
	}

	return table.NewQueryParameters(params...)
}

func MakeWhereClause(filters []YdbFilter) string {
	if len(filters) == 0 {
		return ""
	}

	b := strings.Builder{}

	b.WriteString(" WHERE\n ")

	groupedFilters := groupFiltersByColumnName(filters)

	for _, filterList := range groupedFilters {
		for i, filter := range filterList {
			b.WriteString(filter.YqlColumnName)
			b.WriteString(" ")
			b.WriteString(string(filter.Condition))
			b.WriteString(" ")
			b.WriteString(MakeColumnFilterName(filter.YqlColumnName, strconv.Itoa(i+1)))

			b.WriteString(" ")
			b.WriteString("AND")
			b.WriteString(" ")
		}
	}

	str := b.String()
	trimmedStr := str[:len(str)-5] //Removing the last ' AND ' section of the string
	str += ";"

	return trimmedStr
}

func AddYqlVarDeclarations(filters []YdbFilter) string {
	if len(filters) == 0 {
		return ""
	}

	b := strings.Builder{}

	groupedFilters := groupFiltersByColumnName(filters)

	for _, filters := range groupedFilters {
		for i, filter := range filters {
			yqlVarName := MakeColumnFilterName(filter.YqlColumnName, strconv.Itoa(i+1))
			b.WriteString(Declare(yqlVarName, filter.ConditionValue))
		}
	}

	return b.String()
}
