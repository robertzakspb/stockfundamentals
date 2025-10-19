package orderexec

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Execution struct {
	Id          uuid.UUID
	AccountId   uuid.UUID
	SecurityId  uuid.UUID
	Timestamp   time.Time
	Quantity    float64
	Price       float64
	Description string
}

func NewExecution(accountId uuid.UUID, securityId uuid.UUID, timestamp time.Time, quantity float64, price float64, description string) (Execution, error) {
	if quantity < 0 {
		return Execution{}, fmt.Errorf("quantity is smaller than 0: %f", quantity)
	}
	if price < 0 {
		return Execution{}, fmt.Errorf("price is smaller than 0: %f", price)
	}

	if len(description) > 1000 {
		return Execution{}, fmt.Errorf("description cannot contain more than 1000 characters")
	}

	return Execution{
		Id:          uuid.New(),
		AccountId:   accountId,
		SecurityId:  securityId,
		Timestamp:   timestamp,
		Quantity:    quantity,
		Price:       price,
		Description: description,
	}, nil
}

func (exec *Execution) IsBuyOrder() bool {
	return exec.Quantity > 0
}

func (exec *Execution) IsSellOrder() bool {
	return exec.Quantity < 0
}
