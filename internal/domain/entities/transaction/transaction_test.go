package transaction

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_NewExec_Valid(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "testDesc", "BUY"
	transactionType := "DEPOSIT"

	exec, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertEqual(t, accountId, exec.AccountId)
	test.AssertEqual(t, securityId, exec.SecurityId)
	test.AssertEqual(t, timestamp, exec.Timestamp)
	test.AssertEqual(t, quantity, exec.Quantity)
	test.AssertEqual(t, price, exec.Amount)
	test.AssertEqual(t, description, exec.Description)
	test.AssertEqual(t, side, OrderSideStringValue[exec.Side])
}

func Test_NewExec_InvalidQuantity(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := -10.0, 25.0
	description, side := "testDesc", "BUY"
	transactionType := "DEPOSIT"

	_, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertError(t, err)
}

func Test_NewExec_InvalidPrice(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, -25.0
	description, side := "testDesc", "BUY"
	transactionType := "DEPOSIT"

	_, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertError(t, err)
}

func Test_NewExec_InvalidDescription(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := string(make([]byte, 10001)), "BUY"
	transactionType := "DEPOSIT"

	_, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertError(t, err)
}

func Test_NewExec_InvalidOrderSide(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "test"
	transactionType := "DEPOSIT"

	_, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertError(t, err)
}

func Test_IsBuyOrder_Positive(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "BUY"
	transactionType := "DEPOSIT"

	exec, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertTrue(t, exec.IsBuyOrder())
}

func Test_IsBuyOrder_Negative(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "SELL"
	transactionType := "DEPOSIT"

	exec, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertFalse(t, exec.IsBuyOrder())
}

func Test_IsSellOrder_Positive(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "SELL"
	transactionType := "DEPOSIT"

	exec, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertTrue(t, exec.IsSellOrder())
}

func Test_IsSellOrder_Negative(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "BUY"
	transactionType := "DEPOSIT"

	exec, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertFalse(t, exec.IsSellOrder())
}

func Test_TransactionType_Deposit(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "BUY"
	transactionType := "DEPOSIT"

	transaction, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertTrue(t, transaction.IsDeposit())
}

func Test_TransactionType_Withdrawal(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "BUY"
	transactionType := "WITHDRAWAL"

	transaction, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertTrue(t, transaction.IsWithdrawal())
}

func Test_TransactionType_OrderExecution(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "BUY"
	transactionType := "ORDER_EXECUTION"

	transaction, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertEqual(t, OrderExecution, transaction.Type)
}

func Test_IsDepositOrWithdrawal_Deposit(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "BUY"
	transactionType := "DEPOSIT"

	transaction, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertTrue(t, transaction.IsDepositOrWithdrawal())
}

func Test_IsDepositOrWithdrawal_Withdrawal(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "BUY"
	transactionType := "WITHDRAWAL"

	transaction, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side, transactionType)

	test.AssertNoError(t, err)
	test.AssertTrue(t, transaction.IsDepositOrWithdrawal())
}
