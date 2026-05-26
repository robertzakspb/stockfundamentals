package orderexec

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

	exec, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side)

	test.AssertNoError(t, err)
	test.AssertEqual(t, accountId, exec.AccountId)
	test.AssertEqual(t, securityId, exec.SecurityId)
	test.AssertEqual(t, timestamp, exec.Timestamp)
	test.AssertEqual(t, quantity, exec.Quantity)
	test.AssertEqual(t, price, exec.Price)
	test.AssertEqual(t, description, exec.Description)
	test.AssertEqual(t, side, OrderSideStringValue[exec.Side])
}

func Test_NewExec_InvalidQuantity(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := -10.0, 25.0
	description, side := "testDesc", "BUY"

	_, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side)

	test.AssertError(t, err)
}

func Test_NewExec_InvalidPrice(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, -25.0
	description, side := "testDesc", "BUY"

	_, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side)

	test.AssertError(t, err)
}

func Test_NewExec_InvalidDescription(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := string(make([]byte, 10001)), "BUY"

	_, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side)

	test.AssertError(t, err)
}

func Test_NewExec_InvalidOrderSide(t *testing.T) {
	accountId, securityId := uuid.New(), uuid.New()
	timestamp := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	quantity, price := 10.0, 25.0
	description, side := "description", "test"

	_, err := New(accountId, securityId, timestamp, float64(quantity), float64(price), description, side)

	test.AssertError(t, err)
}
