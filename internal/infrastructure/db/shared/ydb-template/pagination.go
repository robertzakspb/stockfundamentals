package ydbtemplate

import (
	"strconv"
	"strings"
)

// Declares and applies a page size to the provided query.
// func addPageSizeToQuery(yql string, pageSize int) string {
// 	if pageSize < 1 {
// 		return yql
// 	}

// 	//Adding the page size to the query
// 	var sb strings.Builder
// 	sb.WriteString("$pageSize = ")
// 	sb.WriteString(strconv.Itoa(int(pageSize)))
// 	sb.WriteString(";\n")

// 	//Adding the variable at the start of the query
// 	yql = sb.String() + yql

// 	//Adding the page size limit at the end of the query
// 	yql += "\nLIMIT $pageSize;"

// 	return yql
// }

// Declares the pageSize variable at the start
func declarePageSizeVariable(sb *strings.Builder, pageSize int) {
	if pageSize < 1 {
		return
	}

	sb.WriteString("$pageSize = ")
	sb.WriteString(strconv.Itoa(int(pageSize)))
	sb.WriteString(";\n")
}

// Limits the number of required rows by the pageSize parameter (for pagination)
func limitQuerySizeByPageSize(sb *strings.Builder, pageSize int) {
	if pageSize < 1 {
		return
	}

	sb.WriteString("\nLIMIT $pageSize;")
}
