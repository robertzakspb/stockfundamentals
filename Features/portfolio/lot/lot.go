package lot

import "fmt"

type Lot struct {
	SecurityID   string
	ISIN         string
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
		lot.ISIN,
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
func (lot Lot) CurrentReturn(quote float64) float64 {
	if lot.OpeningPrice == 0 {
		return 0
	}
	return (quote - lot.OpeningPrice) / lot.OpeningPrice
}

func (lot Lot) MarketValue(quote float64) float64 {
	return lot.Quantity * quote
}
