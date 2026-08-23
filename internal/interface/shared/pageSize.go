package shared

import (
	"errors"
	"strconv"
)

func ParsePageSizeParameter(queryParams map[string][]string) (int, error) {
	pageSizeString, err := GetFromQueryParams("pageSize", queryParams)
	if err != nil {
		return -1, err
	}

	pageSize, err := strconv.Atoi(pageSizeString)
	if err != nil {
		return -1, err
	}

	if pageSize < 1 {
		return -1, errors.New("The page size must be at least 1; however, the following value ws provided: " + pageSizeString)
	}

	return pageSize, nil
}
