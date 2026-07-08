package quote

import "time"

type Quote struct {
	figi      string
	quote     float64
	timestamp time.Time
	currency  string
}

func New(figi, currency string, timestamp time.Time, quote float64) Quote {
	return Quote{
		figi:      figi,
		quote:     quote,
		timestamp: timestamp,
		currency:  currency,
	}
}

func (q Quote) Quote() float64 {
	return q.quote
}

func (q Quote) Figi() string {
	return q.figi
}

func (q Quote) Currency() string {
	return q.currency
}

func (q Quote) Timestamp() time.Time {
	return q.timestamp
}
