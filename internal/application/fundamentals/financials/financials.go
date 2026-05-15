package financialsservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	dbfinancials "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/financials"
)

func GetFinancialMetrics() ([]financials.FinancialMetric, error) {
	dbMetrics, err := dbfinancials.FetchFinancialMetrics()
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
