package forexservice

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_MakeEcbUrl_Positive(t *testing.T) {
	expectedUrl := "https://data-api.ecb.europa.eu/service/data/EXR/D.USD.EUR.SP00.A?startPeriod=2025-09-16&endPeriod=2026-09-16&format=csvdata"
	startDate := time.Date(2025, 9, 16, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC)

	url := makeEcbApiUrl(startDate, endDate)

	test.AssertEqual(t, expectedUrl, url)
}
