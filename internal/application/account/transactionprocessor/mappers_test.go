package transactionprocessor

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_mapTransactionsToDbModel(t *testing.T) {
	id, accountId := uuid.New(), uuid.New()
	figi := "figi1"
	tType := transaction.OrderExecution
	timestamp := time.Now()
	side := transaction.Buy
	quantity := 10.0
	pricePerUnit := 25.5
	currency := "USD"
	description := "test"

	sampleTransaction := transaction.Transaction{
		Id:           id,
		AccountId:    accountId,
		Figi:         figi,
		Type:         tType,
		Timestamp:    timestamp,
		Side:         side,
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
		Description:  description,
	}

	dbModels := mapTransactionsToDbModel([]transaction.Transaction{sampleTransaction})

	test.AssertEqual(t, 1, len(dbModels))
	test.AssertEqual(t, id, dbModels[0].Id)
	test.AssertEqual(t, accountId, dbModels[0].AccountId)
	test.AssertEqual(t, figi, dbModels[0].Figi)
	test.AssertEqual(t, string(tType), dbModels[0].Type)
	test.AssertEqual(t, timestamp, dbModels[0].Timestamp)
	test.AssertEqual(t, string(side), dbModels[0].Side)
	test.AssertEqual(t, quantity, dbModels[0].Quantity)
	test.AssertEqual(t, pricePerUnit, dbModels[0].PricePerUnit)
	test.AssertEqual(t, currency, dbModels[0].Currency)
	test.AssertEqual(t, description, dbModels[0].Description)
}
