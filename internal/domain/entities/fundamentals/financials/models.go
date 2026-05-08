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
	Q1            ReportingPeriod = "Q1"
	Q2            ReportingPeriod = "Q2"
	Q3            ReportingPeriod = "Q3"
	Q4            ReportingPeriod = "Q4"
	CALENDAR_YEAR ReportingPeriod = "CALENDAR_YEAR"
	H1            ReportingPeriod = "H1"
	H2            ReportingPeriod = "H2"
	NINE_MONTHS   ReportingPeriod = "NINE_MONTHS"
)

var ReportingPeriodMap = map[string]ReportingPeriod{
	"Q1":            Q1,
	"Q2":            Q2,
	"Q3":            Q3,
	"Q4":            Q4,
	"CALENDAR_YEAR": CALENDAR_YEAR,
	"H1":            H1,
	"H2":            H2,
	"NINE_MONTHS":   NINE_MONTHS,
}
