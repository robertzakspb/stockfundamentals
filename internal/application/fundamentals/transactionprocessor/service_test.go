package transactionprocessor

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_ProcessOrderExecutions_EmptySlice(t *testing.T) {
	ts := []transaction.Transaction{}

	err := ProcessStockOrderExecutions(ts)

	test.AssertError(t, err)
}

func Test_ProcessOrderExecutions_InvalidTransactionType(t *testing.T) {
	ts := []transaction.Transaction{
		{
			Type: transaction.Deposit,
		},
	}

	err := ProcessStockOrderExecutions(ts)

	test.AssertError(t, err)
}

