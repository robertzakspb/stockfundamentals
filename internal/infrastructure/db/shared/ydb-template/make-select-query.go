package ydbtemplate

import (
	"errors"
	"reflect"
	"slices"
	"strings"

	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	taghelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/tag-helpers"
)

func generateGetQuery[DB any](filters []ydbfilter.YdbFilter, sortBy, sortDirection, tablePath string, pageSize int) (string, error) {
	sb := strings.Builder{}

	//Adding YQL declarations first
	sb.WriteString(ydbfilter.AddYqlVarDeclarations(filters))

	//As pagination relies on the sorting parameter, it must be provided
	needPagination := pageSize > 0
	if needPagination && sortBy == "" {
		return sb.String(), errors.New("The sortBy parameter cannot be an empty string if pagination is required.")
	}

	//Declaring the page size in case pagination is required
	declarePageSizeVariable(&sb, pageSize)

	sb.WriteString("SELECT ")

	columnNames, err := taghelpers.GetEntityTagValues[DB]("sql")
	if err != nil {
		return sb.String(), err
	}
	slices.Sort(columnNames)

	//Ensuring the provided filters' column names are present in T's sql tag values
	for i := range filters {
		if !slices.Contains(columnNames, filters[i].YqlColumnName) {
			var t DB
			return sb.String(), errors.New("Column " + filters[i].YqlColumnName + " in a YDB filter is not present in the tag values of entity " + reflect.TypeOf(t).Name())
		}
	}

	for i := range columnNames {
		//Adding the column name to the query
		sb.WriteString(columnNames[i])
		sb.WriteString(", ")
	}

	sb.WriteString("FROM ")
	sb.WriteString(tablePath)
	sb.WriteString(" ")
	sb.WriteString(ydbfilter.MakeWhereClause(filters))

	addSortingToQuery(&sb, sortBy, sortDirection)

	limitQuerySizeByPageSize(&sb, pageSize)

	return sb.String(), nil
}

// func generatePaginatedGetQuery[DB any](filters []ydbfilter.YdbFilter, sortBy, tablePath string, pageSize int) (string, error) {
// 	baseQuery, err := generateGetQueryNoPagination[DB](filters, tablePath)
// 	if err != nil {
// 		return baseQuery, err
// 	}

// 	//As pagination relies on the sorting parameter, it must be provided
// 	needPagination := pageSize > 0
// 	if needPagination && sortBy == "" {
// 		return baseQuery, errors.New("The sortBy parameter cannot be an empty string if pagination is required.")
// 	}

// 	baseQuery = addSortingToQuery(baseQuery, sortBy)

// 	baseQuery = addPageSizeToQuery(baseQuery, pageSize)

// 	return baseQuery, nil
// }
