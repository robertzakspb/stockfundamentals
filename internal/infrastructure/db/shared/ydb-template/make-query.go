package ydbtemplate

import (
	"errors"
	"reflect"
	"slices"
	"strings"

	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	taghelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/tag-helpers"
)

func generateGetQuery[T any](filters []ydbfilter.YdbFilter, tablePath string) (string, error) {
	sb := strings.Builder{}

	sb.WriteString(ydbfilter.AddYqlVarDeclarations(filters))
	sb.WriteString("SELECT ")

	columnNames, err := taghelpers.GetEntityTagValues[T]("sql")
	if err != nil {
		return sb.String(), err
	}

	//Ensuring the provided filters' column names are present in T's sql tag values
	for i := range filters {
		if !slices.Contains(columnNames, filters[i].YqlColumnName) {
			var t T
			return sb.String(), errors.New("Column " + filters[i].YqlColumnName + " in a YDB filter is not present in the tag values of entity " + reflect.TypeOf(t).Name())
		}
	}

	for i := range columnNames {
		//Adding the column name to the query
		sb.WriteString(columnNames[i] + ", ")
	}

	sb.WriteString("FROM ")
	sb.WriteString(tablePath + " ")
	sb.WriteString(ydbfilter.MakeWhereClause(filters))

	return sb.String(), nil
}
