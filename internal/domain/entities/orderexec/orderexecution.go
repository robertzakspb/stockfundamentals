package orderexec

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Execution struct {
	Id          uuid.UUID
	AccountId   uuid.UUID
	SecurityId  uuid.UUID
	Timestamp   time.Time
	Side        OrderSide
	Quantity    float64
	Price       float64
	Description string
}

type OrderSide string

const (
	Buy  OrderSide = "BUY"
	Sell OrderSide = "SELL"
)

var OrderSideStringValue = map[OrderSide]string{
	Buy:  "BUY",
	Sell: "SELL",
}
var OrderSideLookup = map[string]OrderSide{
	"BUY":  Buy,
	"SELL": Sell,
}

func New(accountId, securityId uuid.UUID, timestamp time.Time, quantity, price float64, description, side string) (Execution, error) {
	if quantity < 0 {
		return Execution{}, fmt.Errorf("quantity is smaller than 0: %f", quantity)
	}
	if price < 0 {
		return Execution{}, fmt.Errorf("price is smaller than 0: %f", price)
	}
	if len(description) > 1000 {
		return Execution{}, fmt.Errorf("description cannot contain more than 1000 characters")
	}
	orderSide, found := OrderSideLookup[side]
	if !found {
		return Execution{}, errors.New("Invalid order side: " + side)
	}

	return Execution{
		Id:          uuid.New(),
		AccountId:   accountId,
		SecurityId:  securityId,
		Timestamp:   timestamp,
		Quantity:    quantity,
		Price:       price,
		Side:        orderSide,
		Description: description,
	}, nil
}

func (exec *Execution) IsBuyOrder() bool {
	return exec.Quantity > 0
}

func (exec *Execution) IsSellOrder() bool {
	return exec.Quantity < 0
}
