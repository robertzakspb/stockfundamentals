package shared

import ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"

func ParseFiltrationPaginationAndSorting[API any](queryParams map[string][]string) (ParsedApiQuery, error) {
	filters, err := ydbfilter.MapQueryFiltersToYdb[API](queryParams)
	if err != nil {
		return ParsedApiQuery{}, err
	}

	sortByParameter, direction, err := ParseSortByParameter[API](queryParams)
	if err != nil {
		return ParsedApiQuery{}, err
	}

	pageSize, err := ParsePageSizeParameter(queryParams)
	if err != nil {
		return ParsedApiQuery{}, err
	}

	parsedQuery := ParsedApiQuery{
		Filters:       filters,
		SortByColumn:  sortByParameter,
		SortDirection: direction,
		PageSize:      pageSize,
	}
	return parsedQuery, nil
}
