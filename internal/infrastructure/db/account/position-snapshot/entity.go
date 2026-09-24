package positionsnapshotdb

import (
	"time"

	"github.com/google/uuid"
)

type StockPositionSnapshotDbModel struct {
	Figi      string    `sql:"figi"`
	Isin      string    `sql:"isin"`
	AccountId uuid.UUID `sql:"account_id"`
	Date      time.Time `sql:"date"`
	Quantity  float64   `sql:"quantity"`
}
