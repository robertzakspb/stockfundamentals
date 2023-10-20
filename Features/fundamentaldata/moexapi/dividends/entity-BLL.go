package moexdividends

import "time"

type Dividends struct {
	Isin        string
	Ticker      string
	CompanyName string
	Dividends   []Dividend
}

type Dividend struct {
	AmountPaid float64
	Currency   string
	Date       time.Time
}
