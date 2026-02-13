package portfoliodb

import (
	"time"

	"github.com/google/uuid"
)

type LotDb struct {
	Id           uuid.UUID `sql:"id"`
	Figi         string    `sql:"figi"`
	Ticker       string    `sql:"ticker"`
	CompanyName  string    `sql:"company_name"`
	AccountID    uuid.UUID `sql:"account_id"`
	CreatedAt    time.Time `sql:"created_at"`
	UpdatedAt    time.Time `sql:"updated_at"`
	Quantity     float64   `sql:"quantity"`
	PricePerUnit float64   `sql:"price_per_unit"`
	Currency     string    `sql:"currency"`
}
