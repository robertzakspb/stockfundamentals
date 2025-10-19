package dividend

import (
	"fmt"

	"github.com/google/uuid"
)

type Payment struct {
	Id         uuid.UUID
	DividendId uuid.UUID
	AccountId  uuid.UUID
	Amount     float64
}

func NewDividendPayment(divId uuid.UUID, accountId uuid.UUID, amount float64) (Payment, error) {
	if amount <= 0 {
		return Payment{}, fmt.Errorf("invalid dividend distribution amount: %f for dividend %s", amount, divId.String())
	}

	return Payment{
		Id:         uuid.New(),
		DividendId: divId,
		AccountId:  accountId,
		Amount:     amount,
	}, nil
}
