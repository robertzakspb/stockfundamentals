package bondportfolioapi

import "time"

type bondPositionLotDto struct {
	Figi                    string    `json:"figi"`
	Isin                    string    `json:"isin"`
	Name                    string    `json:"name"`
	OpeningDate             time.Time `json:"openingDate"`
	ModificationDate        time.Time `json:"modificationDate"`
	AccountId               string    `json:"accountId"`
	Quantity                float64   `json:"quantity"`
	PricePerUnitPercentage  float64   `json:"pricePerUnitPercentage"`
	CurrentYtm              float64   `json:"currentYTM"`
	SimpleYieldToCallOption float64   `json:"yieldToCallOption"`
	MarketValueInRUB        float64   `json:"marketValueInRUB"`
	QuoteInPercentage       float64   `json:"quoteInPercentage"`
	Ytm                     float64   `json:"ytm"`
}

type timeLineItemDto struct {
	Timestamp time.Time `json:"timestamp"`
	EventName string    `json:"eventName"`
	BondName  string    `json:"bondName"`
}
