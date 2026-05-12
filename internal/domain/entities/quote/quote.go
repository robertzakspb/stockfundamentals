package quote

import "time"

type Quote struct {
	Figi      string
	Quote     float64
	Timestamp time.Time
	Currency  string
}
