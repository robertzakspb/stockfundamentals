package dividend

import (
	"time"

	"github.com/google/uuid"
)

type Dividend struct {
	Id                uuid.UUID `sql:"id"`
	Figi              string    `sql:"stock_id"`
	ActualDPS         float64   `sql:"actual_DPS"`
	ExpectedDPS       float64   `sql:"expected_DPS"`
	Currency          string    `sql:"currency"`
	AnnouncementDate  time.Time `sql:"announcement_date"`
	RecordDate        time.Time `sql:"record_date"`
	PayoutDate        time.Time `sql:"payout_date"`
	PaymentPeriod     string    `sql:"payment_period"`
	ManagementComment string    `sql:"management_comment"`
}