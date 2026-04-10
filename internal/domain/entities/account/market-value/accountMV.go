package accountmvdomain

import (
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/compoundinterest"
	"github.com/google/uuid"
)

type AccountMarketValue struct {
	AccountId uuid.UUID
	Date      time.Time
	Currency  string
	EodValue  float64
}

type Return struct {
	AccountId                uuid.UUID `json:"accountId" sql:"account_id"`
	Currency                 string
	AbsoluteReturn           float64
	AbsoluteReturnPercentage float64   //12% would be set as 0.12
	AnnualizedReturn         float64   //12% would be set as 0.12
	StartDate                time.Time `json:"date" sql:"date"`
	EndDate                  time.Time `json:"date" sql:"date"`
}

// Returns the difference between two market values along with the currency
func CalculateAccountReturn(accountId uuid.UUID, startDateMV, endDateMV AccountMarketValue) Return {
	absoluteReturn := endDateMV.EodValue - startDateMV.EodValue

	absoluteReturnPercentage := 0.0
	if startDateMV.EodValue != 0 {
		absoluteReturnPercentage = absoluteReturn / startDateMV.EodValue
	}

	annualizedReturn := compoundinterest.CalcAnnualizedReturn(absoluteReturnPercentage, startDateMV.Date, endDateMV.Date)

	return Return{
		AccountId:                accountId,
		Currency:                 startDateMV.Currency,
		AbsoluteReturn:           absoluteReturn,
		AbsoluteReturnPercentage: absoluteReturnPercentage,
		AnnualizedReturn:         annualizedReturn,
		StartDate:                startDateMV.Date,
		EndDate:                  endDateMV.Date,
	}
}
