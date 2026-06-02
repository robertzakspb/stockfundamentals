package account

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_IsCashNegative_PositiveCase(t *testing.T) {
	account := Account{CashBalance: -100}

	test.AssertTrue(t, account.IsCashNegative())
}

func Test_IsCashNegative_NegativeCase(t *testing.T) {
	account := Account{CashBalance: 100}

	test.AssertFalse(t, account.IsCashNegative())
}

func Test_IsCashNegative_Zero(t *testing.T) {
	account := Account{CashBalance: 0}

	test.AssertFalse(t, account.IsCashNegative())
}

func Test_IsCashPositive_PositiveCase(t *testing.T) {
	account := Account{CashBalance: 100}

	test.AssertTrue(t, account.IsCashPositive())
}

func Test_IsCashPositive_NegativeCase(t *testing.T) {
	account := Account{CashBalance: -100}

	test.AssertFalse(t, account.IsCashPositive())
}

func Test_IsCashPositive_Zero(t *testing.T) {
	account := Account{CashBalance: 0}

	test.AssertTrue(t, account.IsCashPositive())
}

