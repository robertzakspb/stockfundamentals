package dividend

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Payout struct {
	Id         uuid.UUID
	Figi       string
	Ticker     string
	DividendId uuid.UUID
	AccountId  uuid.UUID
	Amount     float64
	Date       time.Time
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
