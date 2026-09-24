package forexservice

import (
	"encoding/csv"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	typeconverter "github.com/compoundinvest/stockfundamentals/internal/utilities/converters"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
)

// Fetches the /EUR currency exchange rate
func FetchUsdToEurRate(startDate, endDate time.Time) ([]ForexRate, error) {
	if startDate.After(endDate) {
		return []ForexRate{}, errors.New("The provided start date is after the end date, abandoning the USD/EUR fetching process")
	}
	url := makeEcbApiUrl(startDate, endDate)

	res, err := http.Get(url)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}

	reader := csv.NewReader(res.Body)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}

	rates := []ForexRate{}
	for i := range records {
		if i == 0 {
			continue //The first row is headers
		}
		if len(records[i]) < 8 {
			logger.Log("The row unexpectedly contains less than 8 columns and is likely corrupt", logger.ERROR)
			continue
		}
		dateString := records[i][6]
		date, err := timehelpers.DateFromISOstring(dateString)
		if err != nil {
			logger.LogError(err, logger.ERROR)
			continue
		}
		rateString := records[i][7]
		rate, err := typeconverter.GetFloat(rateString)
		if err != nil {
			logger.LogError(err, logger.ERROR)
			continue
		}
		if rate == 0.0 {
			logger.Log("The EUR/USD rate is 0, skipping the row", logger.ERROR)
			continue
		}

		//As the provided rate is EUR/USD, we need to inverse it to find USD/EUR
		inversedRate := 1 / rate

		fxRate := ForexRate{
			Currency1: "USD",
			Currency2: "EUR",
			Rate:      inversedRate,
			Date:      date,
		}
		rates = append(rates, fxRate)
	}

	return rates, nil
}

// "https://data-api.ecb.europa.eu/service/data/EXR/D.USD.EUR.SP00.A?startPeriod=2025-09-16&endPeriod=2026-09-16&format=csvdata"
func makeEcbApiUrl(startDate, endDate time.Time) string {
	var sb strings.Builder
	sb.WriteString("https://data-api.ecb.europa.eu/service/data/EXR/D.USD.EUR.SP00.A?startPeriod=")
	sb.WriteString(timehelpers.DateInIsoFormat(startDate))
	sb.WriteString("&endPeriod=")
	sb.WriteString(timehelpers.DateInIsoFormat(endDate))
	sb.WriteString("&format=csvdata")
	return sb.String()
}
