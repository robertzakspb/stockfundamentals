package appdividend

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
)

func GetFilteredDividends(filters []ydbfilter.YdbFilter) ([]dividend.Dividend, error) {
	return dbdividend.GetFilteredDividends(filters)
}

func GetAllUpcomingDividends() ([]dividend.Dividend, error) {
	return dbdividend.GetUpcomingDividends()
}
