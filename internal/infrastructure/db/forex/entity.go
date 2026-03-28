package forexdb

import "time"

type ForexRateDb struct {
	Currency1 string    `sql:"currency_1"`
	Currency2 string    `sql:"currency_2"`
	Date      time.Time `sql:"date"`
	Rate      float64   `sql:"rate"`
}
