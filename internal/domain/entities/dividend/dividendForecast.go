package dividend

type DividendForecast struct {
	Figi          string  `json:"figi"`
	ExpectedDPS   float64 `json:"expectedDPS"`
	Currency      string  `json:"currency"`
	PaymentPeriod string  `json:"paymentPeriod"`
	Author        string  `json:"forecastAuthor"`
	Comment       string  `json:"comment"`
}
