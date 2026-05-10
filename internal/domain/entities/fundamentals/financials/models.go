package financials

import "github.com/google/uuid"

type FinancialMetric struct {
	Id              uuid.UUID
	StockId         uuid.UUID
	Name            string
	ReportingPeriod Period
	Year            int
	Value           int
	Currency        string
}

type Period string

const (
	Q1            Period = "Q1"
	Q2            Period = "Q2"
	Q3            Period = "Q3"
	Q4            Period = "Q4"
	CALENDAR_YEAR Period = "CALENDAR_YEAR"
	H1            Period = "H1"
	H2            Period = "H2"
	NINE_MONTHS   Period = "NINE_MONTHS"
)

var ReportingPeriodMap = map[string]Period{
	"Q1":            Q1,
	"Q2":            Q2,
	"Q3":            Q3,
	"Q4":            Q4,
	"CALENDAR_YEAR": CALENDAR_YEAR,
	"H1":            H1,
	"H2":            H2,
	"NINE_MONTHS":   NINE_MONTHS,
}
