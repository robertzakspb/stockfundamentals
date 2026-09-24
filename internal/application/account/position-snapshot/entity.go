package positionsnapshot

import (
	"time"

	"github.com/google/uuid"
)

type StockPositionSnapshot struct {
	Figi      string
	Isin      string
	AccountId uuid.UUID
	Date      time.Time
	Quantity  float64
}
