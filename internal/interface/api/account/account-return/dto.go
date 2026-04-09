package accountreturnapi

import (
	"time"

	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/dividend-calendar/market-value"
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

func mapDomainToDto(domain accountmvdomain.Return) AccountReturnDto {
	return AccountReturnDto{
		AccountId:                domain.AccountId.String(),
		Currency:                 domain.Currency,
		AbsoluteReturn:           domain.AbsoluteReturn,
		AbsoluteReturnPercentage: domain.AbsoluteReturnPercentage,
		AnnualizedReturn:         domain.AnnualizedReturn,
		StartDate:                domain.StartDate,
		EndDate:                  domain.EndDate,
	}
}
