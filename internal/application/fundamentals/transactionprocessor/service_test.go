package transactionprocessor

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_ProcessOrderExecutions_EmptySlice(t *testing.T) {
	ts := []transaction.Transaction{}

	err := ProcessOrderExecutions(ts)

	test.AssertError(t, err)
}

func Test_ProcessOrderExecutions_InvalidTransactionType(t *testing.T) {
	ts := []transaction.Transaction{
		{
			Type: transaction.Deposit,
		},
	}

	err := ProcessOrderExecutions(ts)

	test.AssertError(t, err)
}

func Test_ExtractAccountsFrom(t *testing.T) {
	uuid1, uuid2 := uuid.New(), uuid.New()
	ts := []transaction.Transaction{
		{
			Type:      transaction.Deposit,
			AccountId: uuid1,
		},
		{
			Type:      transaction.Deposit,
			AccountId: uuid2,
		},
	}

	accountIds := ExtractAccountsFrom(ts)

	test.AssertEqual(t, 2, len(accountIds))
	test.AssertEqual(t, uuid1, accountIds[0])
	test.AssertEqual(t, uuid2, accountIds[1])
}

func Test_GroupByAccount(t *testing.T) {
	uuid1, uuid2 := uuid.New(), uuid.New()
	transactions := []transaction.Transaction{
		{
			Type:      transaction.Deposit,
			Amount:    10,
			AccountId: uuid1,
		},
		{
			Type:      transaction.Deposit,
			AccountId: uuid2,
			Amount:    25,
		},
		{
			Type:      transaction.Deposit,
			AccountId: uuid1,
			Amount:    50,
		},
	}

	grouped := GroupByAccount(transactions)

	test.AssertEqual(t, 2, len(grouped))
	test.AssertEqual(t, 10, grouped[uuid1][0].Amount)
	test.AssertEqual(t, 25, grouped[uuid2][0].Amount)
	test.AssertEqual(t, 50, grouped[uuid1][1].Amount)
}
