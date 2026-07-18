package accountreturnapi

import (
	"time"

	"github.com/google/uuid"
)

type AccountReturnDto struct {
	AccountId                string    `json:"accountId"`
	Currency                 string    `json:"currency"`
	AbsoluteReturn           float64   `json:"absoluteReturn"`
	AbsoluteReturnPercentage float64   `json:"absoluteReturnPercentage"` //12% would be set as 0.12
	AnnualizedReturn         float64   `json:"annualizedReturn"`         //12% would be set as 0.12
	StartDate                time.Time `json:"startDate"`
	EndDate                  time.Time `json:"endDate"`
}

type AccountMarketValueDto struct {
	AccountId uuid.UUID `json:"accountId"`
	Currency  string    `json:"currency"`
	Date      time.Time `json:"date"`
	EodValue  float64   `json:"eodValue"`
}
