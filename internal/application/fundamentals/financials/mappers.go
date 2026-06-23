package financialsservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	dbfinancials "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/financials"
)

func MapFinancialMetricsModelToDbModels(metrics []financials.FinancialMetric) []dbfinancials.FinancialMetricDbModel {
	dbModels := []dbfinancials.FinancialMetricDbModel{}
	for _, metric := range metrics {
		dbModel := dbfinancials.FinancialMetricDbModel{
			Id:              metric.Id,
			StockId:         metric.StockId,
			Name:            metric.Name,
			ReportingPeriod: string(metric.ReportingPeriod),
			Year:            int64(metric.Year),
			Value:           int64(metric.Value),
			Currency:        metric.Currency,
		}
		dbModels = append(dbModels, dbModel)
	}

	return dbModels
}

func mapYdbMetricsToMetrics(dbMetrics []dbfinancials.FinancialMetricDbModel) []financials.FinancialMetric {
	metrics := make([]financials.FinancialMetric, len(dbMetrics))

	for i, dbMetric := range dbMetrics {
		name, _ := financials.MetricMap[dbMetric.Name]
		period, _ := financials.ReportingPeriodMap[dbMetric.ReportingPeriod]
		metric := financials.FinancialMetric{
			Id:              dbMetric.Id,
			StockId:         dbMetric.StockId,
			Name:            name,
			ReportingPeriod: period,
			Year:            int(dbMetric.Year),
			Value:           int(dbMetric.Value),
			Currency:        dbMetric.Currency,
		}
		metrics[i] = metric
	}

	return metrics
}
