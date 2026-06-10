package transactionsapi

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/google/uuid"
)

func Test_mapTransactionsToDto(t *testing.T) {
	id, accountId := uuid.New(), uuid.New()
	figi := "figi1"
	tType := transaction.OrderExecution
	timestamp := time.Date(2026, 1, 1, 1, 1, 1, 1, time.UTC)
	side := transaction.Buy
	quantity := 15.0
	pricePerUnit := 25.2
	currency := "EUR"
	description := "desc"

	t := transaction.Transaction {
		
	}
}
