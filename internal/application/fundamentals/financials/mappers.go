package financialsservice

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	dbfinancials "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/financials"
)

func MapFinancialMetricModelToDbModel(metrics []financials.FinancialMetric) []dbfinancials.FinancialMetricDbModel {
	dbModels := []dbfinancials.FinancialMetricDbModel{}
	for _, metric := range metrics {
		dbModel := dbfinancials.FinancialMetricDbModel{
			Id:       metric.Id,
			StockId:  metric.StockId,
			Name:     metric.Name,
			Period:   string(metric.ReportingPeriod),
			Year:     int64(metric.Year),
			Value:    int64(metric.Value),
			Currency: metric.Currency,
		}
		dbModels = append(dbModels, dbModel)
	}

	return dbModels
}

func mapYdbMetricToMetric(dbMetric dbfinancials.FinancialMetricDbModel) financials.FinancialMetric {
	return financials.FinancialMetric{
		Id:              dbMetric.Id,
		StockId:         dbMetric.StockId,
		Name:            dbMetric.Name,
		ReportingPeriod: financials.ReportingPeriodMap[dbMetric.Period],
		Year:            int(dbMetric.Year),
		Value:           int(dbMetric.Value),
		Currency:        dbMetric.Currency,
	}
}
