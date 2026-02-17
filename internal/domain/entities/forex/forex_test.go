package forex

import (
	"slices"
	"strings"
	"testing"
)

var mockExchangeRateUsdTo = map[Currency]float64{
	RUB: 80,
	EUR: 0.86,
	RSD: 99.54,
	USD: 1.0,
}

type ForexMockDataProvider struct {
}

// Implementation of the ForexDataProvider interface
func (f ForexMockDataProvider) GetExchangeRateUsdTo(currency string) (float64, error) {
	return mockExchangeRateUsdTo[Currency(currency)], nil
}

func (f ForexMockDataProvider) IsSupportedCurrency(cur string) bool {
	supportedCurrencies := []string{"USD", "RSD", "EUR", "RUB"}
	return slices.Contains(supportedCurrencies, strings.ToUpper(cur))
}

func Test_ExchangeRateForNonUsdCurrencies(t *testing.T) {
	expected := mockExchangeRateUsdTo[Currency("USD")] / mockExchangeRateUsdTo[Currency("RSD")]
	actual, err := GetExchangeRateForPair("USD", "RSD", ForexMockDataProvider{})

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	if expected != actual {
		t.Errorf("Expected: %f, actual: %f", expected, actual)
	}
}

func Test_ConvertPriceToAnotherCurrency(t *testing.T) {
	usdToEur := mockExchangeRateUsdTo[Currency("USD")] / mockExchangeRateUsdTo[Currency("EUR")]
	priceInUsd := 200.0

	expectedPriceInEur := priceInUsd * usdToEur
	actual, err := ConvertPriceToDifferentCurrency(priceInUsd, "USD", "EUR", ForexMockDataProvider{})

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	if expectedPriceInEur != actual {
		t.Errorf("Expected: %f, actual: %f", expectedPriceInEur, actual)
	}
}

func Test_ExchangeRateUsdToUsdIsOne(t *testing.T) {
	mock := ForexMockDataProvider{}
	expected := 1.0
	actual, err := mock.GetExchangeRateUsdTo("USD")

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	if actual != expected {
		t.Errorf("Expected: %f, actual: %f", expected, actual)
	}
}

func Test_ExcahgeRateUsdToEur(t *testing.T) {
	mock := ForexMockDataProvider{}
	expected := 0.86
	actual, err := mock.GetExchangeRateUsdTo("EUR")

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}

	if actual != expected {
		t.Errorf("Expected: %f, actual: %f", expected, actual)
	}

}

func Test_UnsupportedCurrencies(t *testing.T) {
	mock := ForexMockDataProvider{}
	actual := mock.IsSupportedCurrency("CHMF")
	expected := false

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func Test_RubIsSupported(t *testing.T) {
	mock := ForexMockDataProvider{}
	actual := mock.IsSupportedCurrency("RUB")
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func Test_UsdIsSupported(t *testing.T) {
	mock := ForexMockDataProvider{}
	actual := mock.IsSupportedCurrency("USD")
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func Test_EurIsSupported(t *testing.T) {
	mock := ForexMockDataProvider{}
	actual := mock.IsSupportedCurrency("Eur")
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func Test_RsdIsSupported(t *testing.T) {
	mock := ForexMockDataProvider{}
	actual := mock.IsSupportedCurrency("RSD")
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}
