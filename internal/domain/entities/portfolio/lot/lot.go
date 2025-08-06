package lot

import (
	"fmt"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/forex"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/google/uuid"
)

type Lot struct {
	Id           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Quantity     float64 `json:"quantity"`
	PricePerUnit float64 `json:"PricePerUnit"`
	Currency     string  `json:"currency"`
	AccountId    uuid.UUID
	// CompanyName  string `json:"companyName"`
	Security security.Security
	// Figi         string
	// Ticker       string `json:"ticker"`
	// MIC          string `json:"MIC"`
}

func NewLot(security security.Security, quantity float64, pricePerUnit float64, currency string, accountId uuid.UUID) (Lot, error) {
	newLot := Lot{
		Id:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
		AccountId:    accountId,
		Security:     security,
	}

	if err := newLot.validate(); err != nil {
		return newLot, err
	}

	return newLot, nil
}

func (lot *Lot) validate() error {
	if lot.Quantity < 0 {
		return fmt.Errorf("Position with ID %v has an unexpected quantity of %v", lot.Id, lot.Quantity)
	}
	if lot.PricePerUnit < 0 {
		return fmt.Errorf("Position with ID %v has an unexpected price per unit of %v", lot.Id, lot.PricePerUnit)
	}
	if !forex.IsSupportedCurrency(lot.Currency) {
		return fmt.Errorf("Position with ID %v has an unsupported currency", lot.Currency)
	}
	if lot.Security == nil {
		return fmt.Errorf("Position with ID does not have a corresponding security")
	}

	return nil
}

func (lot Lot) CostBasis() float64 {
	return lot.Quantity * lot.PricePerUnit
}

func (lot Lot) MergeWith(newLot Lot) (Lot, error) {
	if lot.Security.GetFigi() != newLot.Security.GetFigi() {
		return Lot{}, fmt.Errorf("attempting to merge two lots with a different underlying security")
	}

	newQuantity := lot.Quantity + newLot.Quantity
	newOpeningPrice := (lot.Quantity*lot.PricePerUnit + newLot.Quantity*newLot.PricePerUnit) / newQuantity

	validatedLot, err := NewLot(lot.Security, newQuantity, newOpeningPrice, lot.Currency, lot.AccountId)

	return validatedLot, err
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
		logger.Log("Quote is nil for position "+lot.Security.GetFigi(), logger.ERROR)
	}

	const targetCur = "EUR"
	quoteInTargerCurrency, err := forex.ConvertPriceToDifferentCurrency(quote.Quote(), quote.Currency(), targetCur)
	if err != nil {
		return 0, err
	}

	return lot.Quantity * quoteInTargerCurrency, nil
}
