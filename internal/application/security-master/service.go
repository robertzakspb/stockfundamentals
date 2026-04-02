package security_master

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	dbsecurity "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/security"
	securitydb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/security"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
)

func GetFilteredSecurities(filters []ydbfilter.YdbFilter) ([]security.Stock, error) {
	return securitydb.GetFilteredSecurities(filters)
}

func GetAllSecuritiesFromDB() ([]security.Stock, error) {
	return securitydb.GetAllSecuritiesFromDB()
}

func GetSecuritiesFilteredByFigi(figis []string) ([]security.Stock, error) {
	return securitydb.GetSecuritiesFilteredByFigi(figis)
}

func GetSecuritiesByIsin(isins []string) ([]security.Stock, error) {
	ydbFigis := ydbhelper.ConvertStringsToYdbList(isins)
	filters := []ydbfilter.YdbFilter{{
		YqlColumnName:  "isin",
		Condition:      ydbfilter.Contains,
		ConditionValue: ydbFigis,
	}}

	stocks, err := dbsecurity.GetFilteredSecurities(filters)
	if err != nil {
		return []security.Stock{}, err
	}

	return stocks, nil
}
