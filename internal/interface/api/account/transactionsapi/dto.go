package transactionsapi

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDto struct {
	Id           uuid.UUID
	AccountId    uuid.UUID
	Figi         string
	Type         string
	Timestamp    time.Time
	Side         string
	Quantity     float64
	PricePerUnit float64
	Currency     string
	Description  string
}
