package entity

import "github.com/google/uuid"

type FinancialMetric struct {
	Id       uuid.UUID
	StockId  uuid.UUID
	Name     string
	Period   ReportingPeriod
	Year     int
	Value    int
	Currency string
}

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

var ReportingPeriodMap = map[string]ReportingPeriod{
	"q1":  q1,
	"q2":  q2,
	"q3":  q3,
	"q4":  q4,
	"fy":  fullYear,
	"h1":  h1,
	"h2":  h2,
	"9mo": nine_months,
}
