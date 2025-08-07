package transaction

import "github.com/google/uuid"

type Transaction struct {
	Id          uuid.UUID
	AccountId   uuid.UUID
	Type        TransactionType
	Quantity    float64
	Price       float64
	Description string
}

type TransactionType string

const (
	OrderExecution TransactionType = "OrderExecution"
	Dividend       TransactionType = "Dividend"
	Deposit        TransactionType = "Deposit"
	Withdrawal     TransactionType = "Withdrawal"
)

var TransactionTypeMap = map[string]TransactionType{
	"OrderExecution": OrderExecution,
	"Dividend":       Dividend,
	"Deposit":        Deposit,
	"Withdrawal":     Withdrawal,
}

func (t *Transaction) isBuyOrder() bool {
	return t.Quantity > 0
}

func (t *Transaction) isSellOrder() bool {
	return t.Quantity < 0
}
