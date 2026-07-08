package timeseriesdb

import "time"

type QuoteDB struct {
	Figi       string    `sql:"figi"`
	Date       time.Time `sql:"date"`
	ClosePrice float64   `sql:"close_price"`
	Country    string    `sql:"country_iso2"`
}

type BondQuoteDB struct {
	Ticker            string    `sql:"ticker"`
	Timestamp         time.Time `sql:"timestamp"`
	PriceAsPercentage float64   `sql:"price_as_percentage"`
}
