package forexservice

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

type FxRate struct {
	Date  string         `xml:"Date,attr"`
	Rates []*CbrCurrency `xml:"Valute"`
}

type CbrCurrency struct {
	ID       string  `xml:"ID,attr"`
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	Nominal  float64 `xml:"Nominal"`
	Name     string  `xml:"Name"`
	Value    string  `xml:"Value"`
}

func getCurrencyToRubRate(currency string, targetDate time.Time) (float64, error) {
	rates, err := getDailyRates(context.TODO(), targetDate.Year(), targetDate.Month(), targetDate.Day())
	if err != nil {
		return -1, err
	}

	for _, rate := range rates.Rates {
		if rate.CharCode == currency {
			parsedValue, err := strconv.ParseFloat(strings.Join(strings.Split(rate.Value, ","), "."), 64)
			if err != nil {
				return -1, err
			}
			return parsedValue, nil
		}
	}

	return -1, errors.New("Failed to find the target forex rate")
}

func getDailyRates(ctx context.Context, year int, month time.Month, day int) (*FxRate, error) {
	const (
		UserAgent = "cbrates/v0 (+https://github.com/robertzak)" // by some reason default go ua was getting blocked
	)

	c := http.Client{}

	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("02.01.2006")
	url := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", UserAgent)

	resp, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result FxRate
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Rates) == 0 {
		return nil, errors.New("no data for this date")
	}

	return &result, nil
}
