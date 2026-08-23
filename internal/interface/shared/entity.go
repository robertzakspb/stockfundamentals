package shared

import ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"

// Ancillary struct for storing the DB filters, sorting, and pagination parsed from API requests
type ParsedApiQuery struct {
	Filters       []ydbfilter.YdbFilter
	SortByColumn  string
	SortDirection string
	PageSize      int
}
