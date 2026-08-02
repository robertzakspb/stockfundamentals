package ydbtemplate

import (
	"strings"

	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
)

func generateDeleteQuery(filters []ydbfilter.YdbFilter, tablePath string) (string) {
	sb := strings.Builder{}

	sb.WriteString(ydbfilter.AddYqlVarDeclarations(filters))
	sb.WriteString("DELETE FROM ")

	sb.WriteString(tablePath)
	sb.WriteString(" ")

	sb.WriteString(ydbfilter.MakeWhereClause(filters))
	return sb.String()
}
