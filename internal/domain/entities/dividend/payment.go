package dividend

import (
	"fmt"

	"github.com/google/uuid"
)

type Payout struct {
	Id         uuid.UUID
	DividendId uuid.UUID
	AccountId  uuid.UUID
	Amount     float64
	Dividend   Dividend
}

func NewDividendPayment(divId uuid.UUID, accountId uuid.UUID, amount float64) (Payout, error) {
	if amount <= 0 {
		return Payout{}, fmt.Errorf("invalid dividend distribution amount: %f for dividend %s", amount, divId.String())
	}

	return Payout{
		Id:         uuid.New(),
		DividendId: divId,
		AccountId:  accountId,
		Amount:     amount,
	}, nil
}
