package lot

import (
	"fmt"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/features/marketdata/forex"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/google/uuid"
)

type Lot struct {
	Id           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	SecurityID   string
	Figi         string
	Ticker       string  `json:"ticker"`
	Quantity     float64 `json:"quantity"`
	PricePerUnit float64 `json:"PricePerUnit"`
	Currency     string  `json:"currency"`
	CompanyName  string  `json:"companyName"`
	BrokerName   string  `json:"brokerName"`
	MIC          string  `json:"MIC"`
}

func (lot Lot) MergeWith(newLot Lot) (Lot, error) {
	if lot.Ticker != newLot.Ticker {
		return Lot{}, fmt.Errorf("attempting to merge two lots with a different underlying security")
	}

	newQuantity := lot.Quantity + newLot.Quantity
	newOpeningPrice := (lot.Quantity*lot.PricePerUnit + newLot.Quantity*newLot.PricePerUnit) / newQuantity

	return Lot{
		Id:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		SecurityID:   "",
		Figi:         lot.Figi,
		Ticker:       lot.Ticker,
		Quantity:     newQuantity,
		PricePerUnit: newOpeningPrice,
		Currency:     lot.Currency,
		CompanyName:  lot.CompanyName,
		BrokerName:   "Multiple",
		MIC:          lot.MIC,
	}, nil
}

// Returns the current profit on the lot given a quote (expressed as a percentage)
func (lot Lot) CurrentReturn(quote entity.SimpleQuote) float64 {
	if lot.PricePerUnit == 0 {
		return 0
	}
	return (quote.Quote() - lot.PricePerUnit) / lot.PricePerUnit
}

func (lot Lot) MarketValue(quote entity.SimpleQuote) (float64, error) {
	if quote == nil {
		logger.Log("Quote is nil for position "+lot.Figi, logger.ERROR)
	}

	const targetCur = "EUR"
	quoteInTargerCurrency, err := forex.ConvertPriceToDifferentCurrency(quote.Quote(), quote.Currency(), targetCur)
	if err != nil {
		return 0, err
	}

	return lot.Quantity * quoteInTargerCurrency, nil
}

func (lot Lot) CostBasis() float64 {
	return  lot.Quantity * lot.PricePerUnit
}
