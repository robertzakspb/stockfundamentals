package ydbtemplate

import (
	"strings"
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_declarePageSizeVariable_Negative_ZeroPageSize(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("someQuery")
	pageSize := 0

	declarePageSizeVariable(&sb, pageSize)

	//As the page size is invalid, the query is expected to remain unchanged
	test.AssertEqual(t, "someQuery", sb.String())
}

func Test_declarePageSizeVariable_Positive_ValidPageSize(t *testing.T) {
	var sb strings.Builder
	pageSize := 5
	expectedQuery := "$pageSize = 5;\n"

	declarePageSizeVariable(&sb, pageSize)

	//As the page size is invalid, the query is expected to remain unchanged
	test.AssertEqual(t, expectedQuery, sb.String())
}

func Test_limitQuerySizeByPageSize_Negative_ZeroPageSize(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM user")
	pageSize := 0

	limitQuerySizeByPageSize(&sb, pageSize)

	test.AssertEqual(t, "SELECT * FROM user", sb.String())
}

func Test_limitQuerySizeByPageSize_Positive_ValidPageSize(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM user")
	pageSize := 5

	limitQuerySizeByPageSize(&sb, pageSize)

	test.AssertEqual(t, "SELECT * FROM user\nLIMIT $pageSize;", sb.String())
}