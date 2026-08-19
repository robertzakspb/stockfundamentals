package shared

func ParseSortByParameter(queryParams map[string][]string) (string, error) {
	return GetFromQueryParams("sortBy", queryParams)
}
