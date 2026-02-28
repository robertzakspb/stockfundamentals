package apidividend

import (
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
)

type DividendDTO struct {
	Figi              string    `json:"figi"`
	ActualDPS         float64   `json:"actualDPS"`
	ExpectedDPS       float64   `json:"expectedDPS"`
	Currency          string    `json:"currency"`
	AnnouncementDate  time.Time `json:"announcementDate"`
	RecordDate        time.Time `json:"recordDate"`
	PayoutDate        time.Time `json:"payoutDate"`
	PaymentPeriod     string    `json:"paymentPeriod"`
	ManagementComment string    `json:"managementComment"`
}

type DividendForecastDTO struct {
	Figi           string  `json:"figi"`
	ExpectedDPS    float64 `json:"expectedDPS"`
	Currency       string  `json:"currency"`
	PaymentPeriod  string  `json:"paymentPeriod"`
	ForecastAuthor string  `json:"forecastAuthor"`
	Comment        string  `json:"comment"`
}

func convertDividendToDTO(dividends []dividend.Dividend) []DividendDTO {
	dtos := []DividendDTO{}
	for _, dividend := range dividends {
		dto := DividendDTO{
			Figi:              dividend.Id.String(),
			ActualDPS:         dividend.ActualDPS,
			ExpectedDPS:       dividend.ExpectedDPS,
			Currency:          dividend.Currency,
			AnnouncementDate:  dividend.AnnouncementDate,
			RecordDate:        dividend.RecordDate,
			PayoutDate:        dividend.PayoutDate,
			PaymentPeriod:     dividend.PaymentPeriod,
			ManagementComment: dividend.ManagementComment,
		}
		dtos = append(dtos, dto)
	}

	return dtos
}
