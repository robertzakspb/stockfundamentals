package shared

import (
	"errors"
	"strings"

	taghelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/tag-helpers"
)

func ParseSortByParameter[API any](queryParams map[string][]string) (column string, direction string, err error) {
	sortByColumnAndDirection, err := GetFromQueryParams("sortBy", queryParams)
	if err != nil {
		return "", "", err
	}

	splitValues := strings.Split(sortByColumnAndDirection, ",")
	if len(splitValues) < 2 {
		return "", "", errors.New("The returned value is missing either the sort by column or the sorting direction")
	}
	direction, column = splitValues[0], splitValues[1]

	sortByColumnSql, err := taghelpers.GetTagValueBySourceTag[API]("json", column, "sql")

	return sortByColumnSql, direction, err
}
