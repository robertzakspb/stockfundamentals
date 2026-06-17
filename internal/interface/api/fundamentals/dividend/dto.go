package apidividend

import (
	"time"
)

type DividendDTO struct {
	Figi              string    `json:"figi"`
	ActualDPS         float64   `json:"actualDPS"`
	ExpectedDPS       float64   `json:"expectedDPS"`
	Currency          string    `json:"currency"`
	AnnouncementDate  time.Time `json:"announcementDate"`
	RecordDate        time.Time `json:"recordDate"`
	PayoutDate        time.Time `json:"payoutDate"`
	PaymentPeriod     string    `json:"paymentPeriod"`
	Regularity        string    `json:"regularity"`
	Type              string    `json:"type"`
	ManagementComment string    `json:"managementComment"`
}

type DividendForecastDTO struct {
	Figi               string    `json:"figi"`
	Ticker             string    `json:"ticker"`
	ExpectedDPS        float64   `json:"expectedDPS"`
	Currency           string    `json:"currency"`
	PaymentPeriod      string    `json:"paymentPeriod"`
	Author             string    `json:"forecastAuthor"`
	Comment            string    `json:"comment"`
	Yield              float64   `json:"yield"`
	ExpectedPayoutDate time.Time `json:"expectedPayoutDate"`
}

type SecurityDivForecastDto struct {
	Figi             string                `json:"figi"`
	Forecasts        []DividendForecastDTO `json:"forecasts"`
	CumulativeReturn float64               `json:"cumulativeReturn"`
}
