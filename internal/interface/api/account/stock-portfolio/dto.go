package portfolio

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioDto struct {
	Lots []LotDto `json:"lots"`
}

type LotDto struct {
	Id            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Quantity      float64   `json:"quantity"`
	PricePerUnit  float64   `json:"pricePerUnit"`
	Currency      string    `json:"currency"`
	AccountId     uuid.UUID `json:"accountId" sql:"account_id"` //ID of the corresponding brokerage account
	Figi          string    `json:"figi"`
	CurrentPL     float64   `json:"currentPL"`
	CurrentReturn float64   `json:"currentReturn"`
	Quote         float64   `json:"quote"`
	Isin          string    `json:"isin"`
	Ticker        string    `json:"ticker"`
	MarketValue   float64   `json:"marketValue"`
}
