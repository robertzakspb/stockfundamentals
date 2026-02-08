package timeseries

import "time"

type QuoteDB struct {
	Figi       string    `sql:"figi"`
	Date       time.Time `sql:"date"`
	Country    string    `sql:"country_iso2"`
}
