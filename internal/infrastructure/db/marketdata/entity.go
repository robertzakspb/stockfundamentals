package timeseriesdb

import "time"

type QuoteDB struct {
	Figi       string    `sql:"figi"`
	Date       time.Time `sql:"date"`
	ClosePrice float64   `sql:"close_price"`
	Country    string    `sql:"country_iso2"`
}

type BondQuoteDB struct {
	Figi              string    `sql:"figi"`
	Ticker            string    `sql:"ticker"`
	Timestamp         time.Time `sql:"timestamp"`
	PriceAsPercentage float64   `sql:"price_as_percentage"`
}

func NewBondQuoteDb(figi, ticker string, timestamp time.Time, priceAsPercentage float64) BondQuoteDB {
	return BondQuoteDB{
		Figi:              figi,
		Ticker:            ticker,
		Timestamp:         timestamp,
		PriceAsPercentage: priceAsPercentage,
	}
}

// Implementation of the entity.BondQuote protocol
func (q BondQuoteDB) GetQuoteAsPercentage() float64 {
	return q.PriceAsPercentage
}

func (q BondQuoteDB) GetYtm() float64 {
	return 0
}

func (q BondQuoteDB) GetTicker() string {
	return q.Ticker
}
func (q BondQuoteDB) GetFigi() string {
	return q.Figi
}
func (q BondQuoteDB) GetTimestamp() time.Time {
	return q.Timestamp
}
