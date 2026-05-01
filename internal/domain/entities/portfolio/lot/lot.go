package lot

import (
	"errors"
	"fmt"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
)

type Lot struct {
	Id           uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Quantity     float64
	PricePerUnit float64
	Currency     string
	AccountId    uuid.UUID `json:"accountId" sql:"account_id"`
	Figi         string
	Quote        float64
	Stock        security.Stock
}

func NewLot(figi string, quantity float64, pricePerUnit float64, currency string, accountId uuid.UUID) (Lot, error) {
	newLot := Lot{
		Id:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		Currency:     currency,
		AccountId:    accountId,
		Figi:         figi,
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
	if !forexservice.IsSupportedCurrency(lot.Currency) {
		return fmt.Errorf("Position with ID %v has an unsupported currency", lot.Currency)
	}
	if lot.Figi == "" {
		return fmt.Errorf("Position with ID does not have a corresponding figi")
	}

	return nil
}

func (lot *Lot) CostBasis() float64 {
	return lot.Quantity * lot.PricePerUnit
}

func (lot *Lot) MergeWith(newLot Lot) (Lot, error) {
	if lot.Figi != newLot.Figi {
		return Lot{}, fmt.Errorf("attempting to merge two lots with a different underlying security")
	}

	newQuantity := lot.Quantity + newLot.Quantity
	newOpeningPrice := (lot.Quantity*lot.PricePerUnit + newLot.Quantity*newLot.PricePerUnit) / newQuantity

	validatedLot, err := NewLot(lot.Figi, newQuantity, newOpeningPrice, lot.Currency, lot.AccountId)

	return validatedLot, err
}

// Returns the current profit on the lot given a quote (expressed as a percentage)
func (lot Lot) CurrentReturn() float64 {
	if lot.PricePerUnit == 0 {
		return 0
	}
	return (lot.Quote - lot.PricePerUnit) / lot.PricePerUnit
}

func (lot Lot) CurrentPL() float64 {
	if lot.PricePerUnit == 0 {
		return 0
	}
	return (lot.Quote - lot.PricePerUnit) * lot.Quantity
}

func (lot Lot) MarketValue() (float64, error) {
	if lot.Quote == 0 {
		logger.Log("Quote is 0 for position "+lot.Figi, logger.ERROR)
		return -1, errors.New("Missing quote for position " + lot.Figi)
	}

	return lot.Quantity * lot.Quote, nil
}
