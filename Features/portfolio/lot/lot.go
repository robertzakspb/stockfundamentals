package lot

import (
	"fmt"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/features/marketdata/forex"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
)

type Lot struct {
	SecurityID   string
	Figi         string
	Ticker       string  `json:"ticker"`
	Quantity     float64 `json:"quantity"`
	OpeningPrice float64 `json:"openingPrice"`
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
	newOpeningPrice := (lot.Quantity*lot.OpeningPrice + newLot.Quantity*newLot.OpeningPrice) / newQuantity

	return Lot{
		"",
		lot.Figi,
		lot.Ticker,
		newQuantity,
		newOpeningPrice,
		lot.Currency,
		lot.CompanyName,
		"mergedPosition",
		lot.MIC,
	}, nil
}

// Returns the current profit on the lot given a quote (expressed as a percentage)
func (lot Lot) CurrentReturn(quote entity.SimpleQuote) float64 {
	if lot.OpeningPrice == 0 {
		return 0
	}
	return (quote.Quote() - lot.OpeningPrice) / lot.OpeningPrice
}

func (lot Lot) MarketValue(quote entity.SimpleQuote) (float64, error) {
	if quote == nil {
		logger.Log("Quote is nil for position " + lot.Figi, logger.ERROR)
	}
	
	const targetCur = "EUR"
	quoteInTargerCurrency, err := forex.ConvertPriceToDifferentCurrency(quote.Quote(), quote.Currency(), targetCur)
	if err != nil {
		return 0, err
	}

	return lot.Quantity * quoteInTargerCurrency, nil
}