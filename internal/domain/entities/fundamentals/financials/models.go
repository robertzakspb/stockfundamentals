package financials

import "github.com/google/uuid"

type FinancialMetric struct {
	Id              uuid.UUID
	StockId         uuid.UUID
	Name            Metric
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

type Metric string

const (
	Revenue      Metric = "REVENUE"
	EBITDA       Metric = "EBITDA"
	OperatingIncome Metric = "OPERATING_INCOME"
	NetIncome    Metric = "NET_INCOME"
	FreeCashFlow Metric = "FREE_CASH_FLOW"
	Dividend     Metric = "DIVIDEND"
	NetDebt      Metric = "NET_DEBT"
)

var MetricMap = map[string]Metric{
	"REVENUE":        Revenue,
	"EBITDA":         EBITDA,
	"NET_INCOME":     NetIncome,
	"FREE_CASH_FLOW": FreeCashFlow,
	"DIVIDEND":       Dividend,
	"NET_DEBT":       NetDebt,
	"OPERATING_INCOME": OperatingIncome,
}
