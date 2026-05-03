package appdividend

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
)

func GetFilteredDividends(filters []ydbfilter.YdbFilter) ([]dividend.Dividend, error) {
	dbDivs, err := dbdividend.GetFilteredDividends(filters)
	if err != nil {
		return []dividend.Dividend{}, err
	}

	mappedDivs := mapDbModelToDividend(dbDivs)

	return mappedDivs, nil
}

func GetAllUpcomingDividends() ([]dividend.Dividend, error) {
	dbDivs, err := dbdividend.GetUpcomingDividends()
	if err != nil {
		return []dividend.Dividend{}, err
	}

	mappedDivs := mapDbModelToDividend(dbDivs)

	return mappedDivs, nil
}