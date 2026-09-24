package bondportfolioapi

import (
	"time"

	"github.com/google/uuid"
)

type bondPositionLotDto struct {
	Figi                   string    `json:"figi" sql:"figi"`
	Isin                   string    `json:"isin" sql:"isin"`
	Name                   string    `json:"name"`
	OpeningDate            time.Time `json:"openingDate" sql:"opening_date"`
	ModificationDate       time.Time `json:"modificationDate" sql:"modification_date"`
	AccountId              uuid.UUID `json:"accountId" sql:"account_id"`
	Quantity               float64   `json:"quantity"`
	PricePerUnitPercentage float64   `json:"pricePerUnitPercentage"`
	MarketValueInRUB       float64   `json:"marketValueInRUB"`
	QuoteInPercentage      float64   `json:"quoteInPercentage"`
	Ytm                    float64   `json:"ytm"`
	YieldToCallOption      float64   `json:"yieldToCallOption"`
	CurrentCouponYield     float64   `json:"currentCouponYield"`
}

type timeLineItemDto struct {
	Timestamp time.Time `json:"timestamp"`
	EventName string    `json:"eventName"`
	BondName  string    `json:"bondName"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
}
