package forex

import "testing"

func Test_ExchangeRateUsdToUsdIsOne(t *testing.T) {
	actual, err := GetExchangeRateUsdTo("USD")
	expected := 1.0

	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}
	
	if actual != expected {
		t.Errorf("Expected: %f, actual: %f", expected, actual)
	}
}

func Test_UnsupportedCurrencies(t *testing.T) {
	actual := IsSupportedCurrency("CHMF") 
	expected := false

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func Test_RubIsSupported(t *testing.T) {
	actual := IsSupportedCurrency("RUB") 
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func Test_UsdIsSupported(t *testing.T) {
	actual := IsSupportedCurrency("USD") 
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func Test_EurIsSupported(t *testing.T) {
	actual := IsSupportedCurrency("Eur") 
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func Test_RsdIsSupported(t *testing.T) {
	actual := IsSupportedCurrency("RSD") 
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}