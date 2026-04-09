package accountmvdb

import (
	"time"

	"github.com/google/uuid"
)

type AccountMarketValueDB struct {
	AccountId uuid.UUID `sql:"account_id"`
	Date      time.Time `sql:"date"`
	Currency  string    `sql:"currency"`
	EodValue  float64   `sql:"eod_value"`
}
