package apifinancials

import (
	"github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/financials"
	entity "github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
)

type FinancialMetricDTO struct {
	Id       string
	StockId  string
	Name     string
	Period   string
	Year     int
	Value    int
	Currency string
}

func mapFinancialMetricToDto(metric entity.FinancialMetric) FinancialMetricDTO {
	return FinancialMetricDTO{
		Id:       metric.Id.String(),
		StockId:  metric.StockId.String(),
		Name:     metric.Name,
		Period:   string(metric.Period),
		Year:     metric.Year,
		Value:    metric.Value,
		Currency: metric.Currency,
	}
}

func FetchFinancialsFromDB() ([]FinancialMetricDTO, error) {
	metrics, err := financials.FetchFinancialMetrics()
	if err != nil {
		return []FinancialMetricDTO{}, err
	}

	dtos := []FinancialMetricDTO{}
	for _, metric := range metrics {
		dtos = append(dtos, mapFinancialMetricToDto(metric))
	}

	return dtos, nil
}
