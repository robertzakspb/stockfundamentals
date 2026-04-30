package dbdividend

import (
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/google/uuid"
)

type DividendDbModel struct {
	Id   uuid.UUID `sql:"id"`
	Figi string    `sql:"stock_id" json:"figi"`
	//For DB-related reasons, expected and actual DPS are converted to integers to remove the fractional part. Multiplied by a million for better accuracy. When reading the value, it must consequently be divided by a million
	ExpectedDpsTimesMillion int64     `sql:"expected_DPS" json:"expectedDPS"`
	ActualDPSTimesMillion   int64     `sql:"actual_DPS" json:"actualDPS"`
	Currency                string    `sql:"currency" json:"currency"`
	AnnouncementDate        time.Time `sql:"announcement_date" json:"announcementDate"`
	RecordDate              time.Time `sql:"record_date" json:"recordDate"`
	PayoutDate              time.Time `sql:"payout_date" json:"payoutDate"`
	PaymentPeriod           string    `sql:"payment_period" json:"paymentPeriod"`
	ManagementComment       string    `sql:"management_comment" json:"managementComment"`
}

func mapDividendToDbModel(dividends []dividend.Dividend) []DividendDbModel {
	dbModels := make([]DividendDbModel, len(dividends))
	for i, dividend := range dividends {
		dbModel := DividendDbModel{
			Id:                      dividend.Id,
			Figi:                    dividend.Figi,
			ExpectedDpsTimesMillion: int64(dividend.ExpectedDPS * 1_000_000),
			ActualDPSTimesMillion:   int64(dividend.ActualDPS * 1_000_000),
			Currency:                dividend.Currency,
			AnnouncementDate:        dividend.AnnouncementDate,
			RecordDate:              dividend.RecordDate,
			PayoutDate:              dividend.PayoutDate,
			PaymentPeriod:           dividend.PaymentPeriod,
			ManagementComment:       dividend.ManagementComment,
		}
		dbModels[i] = dbModel
	}

	return dbModels
}
