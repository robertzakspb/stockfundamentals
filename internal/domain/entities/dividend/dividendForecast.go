package dividend

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/google/uuid"
)

type DividendForecast struct {
	Id            uuid.UUID `json:"id"`
	Stock      security.Stock
	ExpectedDPS   float64 `json:"expectedDPS"`
	Currency      string  `json:"currency"`
	PaymentPeriod string  `json:"paymentPeriod"`
	Author        string  `json:"forecastAuthor"`
	Comment       string  `json:"comment"`
	Yield         float64 `json:"yield"`
}
