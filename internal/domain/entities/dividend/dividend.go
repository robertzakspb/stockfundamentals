package dividend

import (
	"errors"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/forex"
	"github.com/google/uuid"
)

type Dividend struct {
	Id                uuid.UUID
	Figi              string    `sql:"stock_id" json:"figi"`
	ActualDPS         float64   `sql:"actual_DPS" json:"actualDPS"`
	ExpectedDPS       float64   `sql:"expected_DPS" json:"expectedDPS"`
	Currency          string    `sql:"currency" json:"currency"`
	AnnouncementDate  time.Time `sql:"announcement_date" json:"announcementDate"`
	RecordDate        time.Time `sql:"record_date" json:"recordDate"`
	PayoutDate        time.Time `sql:"payout_date" json:"payoutDate"`
	PaymentPeriod     string    `sql:"payment_period" json:"paymentPeriod"` //TODO: Implement or copy from financial reports
	ManagementComment string    `sql:"management_comment" json:"managementComment"`
}

func NewDividend(div Dividend) (Dividend, error) {
	return div, div.validate()
}

func (d *Dividend) validate() error {
	if d.Id == uuid.Nil {
		return errors.New("Missing dividend ID")
	}
	if d.Figi == "" {
		return errors.New("Missing figi")
	}
	if d.ActualDPS < 0 {
		return errors.New("Invalid actual dividend amount")
	}
	if d.ExpectedDPS < 0 {
		return errors.New("Invalid expected dividend amount")
	}
	f := forex.ForexDP{}
	if f.IsSupportedCurrency(d.Currency) == false {
		return errors.New("Unsupported currency")
	}
	if d.AnnouncementDate.Unix() == 0 {
		return errors.New("Invalid announcement date")
	}
	if d.RecordDate.Unix() == 0 {
		return errors.New("Invalid record date")
	}
	if d.PayoutDate.Unix() == 0 {
		return errors.New("Invalid payout date")
	}

	return nil
}
