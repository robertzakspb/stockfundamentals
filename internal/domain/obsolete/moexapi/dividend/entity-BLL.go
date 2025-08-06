package dividend

import "time"

type Dividend struct {
	Isin       string
	Ticker     string
	AmountPaid float64
	Currency   string
	Date       time.Time
}
