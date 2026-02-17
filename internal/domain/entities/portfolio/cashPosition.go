package portfolio

import (
	"errors"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/forex"
	"github.com/google/uuid"
)

type CashPosition struct {
	Id        uuid.UUID
	AccountId uuid.UUID
	Amount    float64
	Currency  string
}

func NewCashPosition(accountId string, amount float64, currency string) (CashPosition, error) {
	parsedAccountId, err := uuid.Parse(accountId)
	if err != nil {
		return CashPosition{}, errors.New("Provided account ID for the initialized cash position is not a valid UUID: " + accountId)
	}

	cashPosition := CashPosition{
		Id:        uuid.New(),
		AccountId: parsedAccountId,
		Amount:    amount,
		Currency:  currency,
	}

	if err := cashPosition.validate(); err != nil {
		return cashPosition, err
	}

	return cashPosition, nil
}

func (cash *CashPosition) validate() error {
	forex := forex.ForexDP{}
	if !forex.IsSupportedCurrency(cash.Currency) {
		return errors.New("Initialized cash position has an unsupported currency: " + cash.Currency)
	}

	return nil
}
