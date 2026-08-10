package forexservice

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_makeNbsBaseUrl_Positive(t *testing.T) {
	date := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	expectedURL := "https://webappcenter.nbs.rs/ExchangeRateWebApp/ExchangeRate/IndexByDate?isSearchExecuted=true&Date=31.12.2025.&ExchangeRateListTypeID=3"
	url := makeNbsUrl(date)

	test.AssertEqual(t, expectedURL, url)
}
