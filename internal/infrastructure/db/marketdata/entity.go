package timeseriesdb

import "time"

type QuoteDB struct {
	Figi       string    `sql:"figi"`
	Date       time.Time `sql:"date"`
	ClosePrice float64   `sql:"close_price"`
	Country    string    `sql:"country_iso2"`
}
