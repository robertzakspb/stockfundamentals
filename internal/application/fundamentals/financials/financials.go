package financialsservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	dbfinancials "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/financials"
	db "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
	ydbtemplate "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-template"
)

func GetFinancialMetrics() ([]financials.FinancialMetric, error) {
	filters := []ydbfilter.YdbFilter{}
	tablePath := ydbhelper.GenerateTablePath(db.STOCK_DIRECTORY_PREFIX, db.FINANCIAL_METRIC_TABLE_NAME)
	dbMetrics, err := ydbtemplate.GetEntity[dbfinancials.FinancialMetricDbModel](filters, tablePath)
	if err != nil {
		return []financials.FinancialMetric{}, err
	}

	mappedMetrics := mapYdbMetricsToMetrics(dbMetrics)

	return mappedMetrics, nil
}

func SaveFinancialMetrics(metrics []financials.FinancialMetric) error {
	dbModels := MapFinancialMetricsModelToDbModels(metrics)

	err := dbfinancials.SaveFinancialMetricsToDb(dbModels)

	return err
}
