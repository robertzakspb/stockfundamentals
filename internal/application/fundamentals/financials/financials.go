package financialsservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	dbfinancials "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/financials"
)

// "github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
// "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/financials"

// func FetchFinancialMetrics() ([]financials.FinancialMetric, error) {
// 	return dbfinancials.FetchFinancialMetrics()
// }

func SaveFinancialMetrics(metrics []financials.FinancialMetric) error {
	dbModels := MapFinancialMetricModelToDbModel(metrics)

	err := dbfinancials.SaveFinancialMetricsToDb(dbModels)

	return err
}
