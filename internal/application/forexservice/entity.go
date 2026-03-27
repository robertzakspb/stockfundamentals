package forexservice

import "time"

type ForexRate struct {
	Currency1 string
	Currency2 string
	Rate      float64
	Date      time.Time
}
