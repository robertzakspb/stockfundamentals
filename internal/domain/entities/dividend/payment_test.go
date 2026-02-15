package dividend

import (
	"testing"

	"github.com/google/uuid"
)

func Test_DividendPaymentCannotBeNegative(t *testing.T) {
	_, err := NewDividendPayment(uuid.Nil, uuid.Nil, -5)
	if err == nil {
		t.Errorf("Expected an error when initializing a dividend payment with a negative value")
	}
}

func Test_DividendPaymentInit(t *testing.T) {
	const divAmount = 10.0
	divId := uuid.New()
	accountId := uuid.New()

	dividend, err := NewDividendPayment(divId, accountId, divAmount)
	if err != nil {
		t.Errorf("Expected an error when initializing a dividend payment with a negative value")
	}

	if dividend.Amount != divAmount {
		t.Errorf("Expected: %f, actual: %f", divAmount, dividend.Amount)
	}
	if divId != dividend.DividendId {
		t.Errorf("Expected: %s, actual: %s", divId.String(), dividend.DividendId.String())
	}
	if dividend.AccountId != accountId {
		t.Errorf("Expected: %s, actual: %s", accountId.String(), dividend.AccountId.String())
	}
}
