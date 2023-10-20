package fundamentals

type StockFundamentals struct {
	Ticker      string
	CompanyName string
	Financials  []StockFinancialResult
}

type StockFinancialResult struct {
	Metric MetricName
	Values []FinancialMetric
}

type MetricName string

const (
	Dividend   MetricName = "dividend"
	Revenue    MetricName = "revenue"
	FCF        MetricName = "fcf"
	EPS        MetricName = "eps"
	DilutedEPS MetricName = "dilutedEPS"
)
