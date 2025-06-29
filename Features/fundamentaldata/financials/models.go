package financials

import "github.com/google/uuid"

type ReportingPeriod string

const (
	q1          ReportingPeriod = "q1"
	q2          ReportingPeriod = "q2"
	q3          ReportingPeriod = "q3"
	q4          ReportingPeriod = "q4"
	fullYear    ReportingPeriod = "fy"
	h1          ReportingPeriod = "h1"
	h2          ReportingPeriod = "h2"
	nine_months ReportingPeriod = "9mo"
)

var (
	ReportingPeriodMap = map[string]ReportingPeriod{
		"q1":  q1,
		"q2":  q2,
		"q3":  q3,
		"q4":  q4,
		"fy":  fullYear,
		"h1":  h1,
		"h2":  h2,
		"9mo": nine_months,
	}
)

type FinancialMetricDbModel struct {
	Id       uuid.UUID `sql:"id"`
	StockId  uuid.UUID `sql:"stock_id"`
	Name     string    `sql:"metric"`
	Period   string    `sql:"reporting_period"`
	Year     int64     `sql:"year"`
	Value    int64     `sql:"metric_value"`
	Currency string    `sql:"metric_currency"`
}

type FinancialMetric struct {
	Id       uuid.UUID
	StockId  uuid.UUID
	Name     string
	Period   ReportingPeriod
	Year     int
	Value    int
	Currency string
}

func mapFinancialMetricModelToDbModel(metrics []FinancialMetric) []FinancialMetricDbModel {
	dbModels := []FinancialMetricDbModel{}
	for _, metric := range metrics {
		dbModel := FinancialMetricDbModel{
			Id:       metric.Id,
			StockId:  metric.StockId,
			Name:     metric.Name,
			Period:   string(metric.Period),
			Year:     int64(metric.Year),
			Value:    int64(metric.Value),
			Currency: metric.Currency,
		}
		dbModels = append(dbModels, dbModel)
	}

	return dbModels
}
