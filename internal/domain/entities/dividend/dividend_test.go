package dividend

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_DividendWithMissingIdIsInvalid(t *testing.T) {
	_, err := NewDividend(Dividend{
		Id:                uuid.Nil,
		Figi:              "abcdef",
		ActualDPS:         100,
		ExpectedDPS:       150,
		Currency:          "USD",
		AnnouncementDate:  time.Now(),
		RecordDate:        time.Now(),
		PayoutDate:        time.Now(),
		PaymentPeriod:     "q1",
		ManagementComment: "",
	})
	expectedError := "Missing dividend ID"

	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, actual: %s", expectedError, err.Error())
	}
}

func Test_DividendWithMissingFigiIsInvalid(t *testing.T) {
	_, err := NewDividend(Dividend{
		Id:                uuid.New(),
		Figi:              "",
		ActualDPS:         100,
		ExpectedDPS:       150,
		Currency:          "USD",
		AnnouncementDate:  time.Now(),
		RecordDate:        time.Now(),
		PayoutDate:        time.Now(),
		PaymentPeriod:     "q1",
		ManagementComment: "",
	})
	expectedError := "Missing figi"

	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, actual: %s", expectedError, err.Error())
	}
}

func Test_DividendWithNegativeActualDpdIsInvalid(t *testing.T) {
	_, err := NewDividend(Dividend{
		Id:                uuid.New(),
		Figi:              "abc",
		ActualDPS:         -100,
		ExpectedDPS:       150,
		Currency:          "USD",
		AnnouncementDate:  time.Now(),
		RecordDate:        time.Now(),
		PayoutDate:        time.Now(),
		PaymentPeriod:     "q1",
		ManagementComment: "",
	})
	expectedError := "Invalid actual dividend amount"

	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, actual: %s", expectedError, err.Error())
	}
}

func Test_DividendWithNegativeExpectedDpdIsInvalid(t *testing.T) {
	_, err := NewDividend(Dividend{
		Id:                uuid.New(),
		Figi:              "abc",
		ActualDPS:         100,
		ExpectedDPS:       -150,
		Currency:          "USD",
		AnnouncementDate:  time.Now(),
		RecordDate:        time.Now(),
		PayoutDate:        time.Now(),
		PaymentPeriod:     "q1",
		ManagementComment: "",
	})
	expectedError := "Invalid expected dividend amount"

	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, actual: %s", expectedError, err.Error())
	}
}

func Test_DividendWitUnsupportedCurrencyIsInvalid(t *testing.T) {
	_, err := NewDividend(Dividend{
		Id:                uuid.New(),
		Figi:              "abc",
		ActualDPS:         100,
		ExpectedDPS:       150,
		Currency:          "CHF",
		AnnouncementDate:  time.Now(),
		RecordDate:        time.Now(),
		PayoutDate:        time.Now(),
		PaymentPeriod:     "q1",
		ManagementComment: "",
	})
	expectedError := "Unsupported currency"

	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, actual: %s", expectedError, err.Error())
	}
}

func Test_DividendWithInvalidAnnouncementDateIsInvalid(t *testing.T) {
	_, err := NewDividend(Dividend{
		Id:                uuid.New(),
		Figi:              "abc",
		ActualDPS:         100,
		ExpectedDPS:       150,
		Currency:          "USD",
		AnnouncementDate:  time.Unix(0,0),
		RecordDate:        time.Now(),
		PayoutDate:        time.Now(),
		PaymentPeriod:     "q1",
		ManagementComment: "",
	})
	expectedError := "Invalid announcement date"

	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, actual: %s", expectedError, err.Error())
	}
}

func Test_DividendWithInvalidRecordDateIsInvalid(t *testing.T) {
	_, err := NewDividend(Dividend{
		Id:                uuid.New(),
		Figi:              "abc",
		ActualDPS:         100,
		ExpectedDPS:       150,
		Currency:          "USD",
		AnnouncementDate:  time.Now(),
		RecordDate:        time.Unix(0,0),
		PayoutDate:        time.Now(),
		PaymentPeriod:     "q1",
		ManagementComment: "",
	})
	expectedError := "Invalid record date"

	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, actual: %s", expectedError, err.Error())
	}
}

func Test_DividendWithInvalidPayoutDateIsInvalid(t *testing.T) {
	_, err := NewDividend(Dividend{
		Id:                uuid.New(),
		Figi:              "abc",
		ActualDPS:         100,
		ExpectedDPS:       150,
		Currency:          "USD",
		AnnouncementDate:  time.Now(),
		RecordDate:        time.Now(),
		PayoutDate:        time.Unix(0,0),
		PaymentPeriod:     "q1",
		ManagementComment: "",
	})
	expectedError := "Invalid payout date"

	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, actual: %s", expectedError, err.Error())
	}
}