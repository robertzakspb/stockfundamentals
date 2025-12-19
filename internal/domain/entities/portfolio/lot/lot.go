package lot

import (
	"fmt"
	"time"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/forex"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
)

type Lot struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Quantity     float64   `json:"quantity"`
	PricePerUnit float64   `json:"pricePerUnit"`
	Currency     string    `json:"currency"`
	AccountId    uuid.UUID `json:"accountId"`
	SecurityId   uuid.UUID `json:"securityId"`
	CurrentPL    float64   `json:"currentPL"`
}

func NewLot(securityId uuid.UUID, quantity float64, pricePerUnit float64, currency string, accountId uuid.UUID) (Lot, error) {
	newLot := Lot{
		Id:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
		AccountId:    accountId,
		SecurityId:   securityId,
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
	if lot.SecurityId == uuid.Nil {
		return fmt.Errorf("Position with ID does not have a corresponding security")
	}

	return nil
}

func (lot Lot) CostBasis() float64 {
	return lot.Quantity * lot.PricePerUnit
}

func (lot Lot) MergeWith(newLot Lot) (Lot, error) {
	if lot.SecurityId.String() != newLot.SecurityId.String() {
		return Lot{}, fmt.Errorf("attempting to merge two lots with a different underlying security")
	}

	newQuantity := lot.Quantity + newLot.Quantity
	newOpeningPrice := (lot.Quantity*lot.PricePerUnit + newLot.Quantity*newLot.PricePerUnit) / newQuantity

	validatedLot, err := NewLot(lot.SecurityId, newQuantity, newOpeningPrice, lot.Currency, lot.AccountId)

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
		logger.Log("Quote is nil for position "+lot.SecurityId.String(), logger.ERROR)
	}

	const targetCur = "EUR"
	quoteInTargerCurrency, err := forex.ConvertPriceToDifferentCurrency(quote.Quote(), quote.Currency(), targetCur)
	if err != nil {
		return 0, err
	}

	return lot.Quantity * quoteInTargerCurrency, nil
}
