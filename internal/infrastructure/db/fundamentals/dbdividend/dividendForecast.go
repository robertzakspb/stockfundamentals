package dbdividend

import "github.com/google/uuid"

type DividendForecastDb struct {
	Id            uuid.UUID `json:"id" sql:"id"`
	Figi          string    `json:"figi" sql:"figi"`
	ExpectedDPS   float64   `json:"expectedDPS" sql:"expected_DPS"`
	Currency      string    `json:"currency" sql:"currency"`
	PaymentPeriod string    `json:"paymentPeriod" sql:"payment_period"`
	Author        string    `json:"forecastAuthor" sql:"forecast_author"`
	Comment       string    `json:"comment" sql:"comment"`
}
