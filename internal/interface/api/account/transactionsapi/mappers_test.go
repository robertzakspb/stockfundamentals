package transactionsapi

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/test"
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

	tran := transaction.Transaction{
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

	dtos := mapTransactionsToDto([]transaction.Transaction{tran})

	test.AssertEqual(t, 1, len(dtos))
	test.AssertEqual(t, id, dtos[0].Id)
	test.AssertEqual(t, accountId, dtos[0].AccountId)
	test.AssertEqual(t, figi, dtos[0].Figi)
	test.AssertEqual(t, string(tType), dtos[0].Type)
	test.AssertEqual(t, string(side), dtos[0].Side)
	test.AssertEqual(t, timestamp, dtos[0].Timestamp)
	test.AssertEqual(t, quantity, dtos[0].Quantity)
	test.AssertEqual(t, pricePerUnit, dtos[0].PricePerUnit)
	test.AssertEqual(t, currency, dtos[0].Currency)
	test.AssertEqual(t, description, dtos[0].Description)
}

func Test_mapDtosToTransactions_Positive(t *testing.T) {
	id, accountId := uuid.New(), uuid.New()
	figi := "figi1"
	tType := transaction.OrderExecution
	timestamp := time.Date(2026, 1, 1, 1, 1, 1, 1, time.UTC)
	side := transaction.Buy
	quantity := 15.0
	pricePerUnit := 25.2
	currency := "EUR"
	description := "desc"

	dto := TransactionDto{
		Id:           id,
		AccountId:    accountId,
		Figi:         figi,
		Type:         string(tType),
		Timestamp:    timestamp,
		Side:         string(side),
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
		Description:  description,
	}

	dtos, err := mapTransactionDtosToTransactions([]TransactionDto{dto})

	test.AssertNoError(t, err)
	test.AssertEqual(t, 1, len(dtos))
	test.AssertEqual(t, id, dtos[0].Id)
	test.AssertEqual(t, accountId, dtos[0].AccountId)
	test.AssertEqual(t, figi, dtos[0].Figi)
	test.AssertEqual(t, tType, dtos[0].Type)
	test.AssertEqual(t, side, dtos[0].Side)
	test.AssertEqual(t, timestamp, dtos[0].Timestamp)
	test.AssertEqual(t, quantity, dtos[0].Quantity)
	test.AssertEqual(t, pricePerUnit, dtos[0].PricePerUnit)
	test.AssertEqual(t, currency, dtos[0].Currency)
	test.AssertEqual(t, description, dtos[0].Description)
}

func Test_mapDtosToTransactions_NonExistentOrderSide(t *testing.T) {
	id, accountId := uuid.New(), uuid.New()
	figi := "figi1"
	tType := transaction.OrderExecution
	timestamp := time.Date(2026, 1, 1, 1, 1, 1, 1, time.UTC)
	quantity := 15.0
	pricePerUnit := 25.2
	currency := "EUR"
	description := "desc"

	dto := TransactionDto{
		Id:           id,
		AccountId:    accountId,
		Figi:         figi,
		Type:         string(tType),
		Timestamp:    timestamp,
		Side:         "FOO",
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
		Description:  description,
	}

	_, err := mapTransactionDtosToTransactions([]TransactionDto{dto})

	test.AssertError(t, err)
}

func Test_mapDtosToTransactions_NonExistentTransactionType(t *testing.T) {
	id, accountId := uuid.New(), uuid.New()
	figi := "figi1"
	timestamp := time.Date(2026, 1, 1, 1, 1, 1, 1, time.UTC)
	quantity := 15.0
	pricePerUnit := 25.2
	currency := "EUR"
	description := "desc"

	dto := TransactionDto{
		Id:           id,
		AccountId:    accountId,
		Figi:         figi,
		Type:         "FOO",
		Timestamp:    timestamp,
		Side:         "BUY",
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
		Description:  description,
	}

	_, err := mapTransactionDtosToTransactions([]TransactionDto{dto})

	test.AssertError(t, err)
}

