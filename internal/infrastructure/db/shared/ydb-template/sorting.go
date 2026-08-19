package ydbtemplate

func addSortingToQuery(yql, sortBy string) string {
	if sortBy == "" {
		return yql
	}

	sortByQuery := "\nORDER BY " + sortBy + "\n"
	yql += sortByQuery

	return yql
}