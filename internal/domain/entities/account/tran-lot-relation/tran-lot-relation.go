package tranlotrelation

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Relationship between a transaction and the bond/stock lots it impacted through purchases and sales
type TransactionLotRelation struct {
	StockLotId    uuid.UUID
	BondLotId     uuid.UUID
	Date          time.Time
	Quantity      float64
}

func New(stockLotId, bondLotId uuid.UUID, date time.Time, quantity float64) (TransactionLotRelation, error) {
	t := TransactionLotRelation{
		StockLotId:    stockLotId,
		BondLotId:     bondLotId,
		Date:          date,
		Quantity:      quantity,
	}

	err := t.validate()

	return t, err
}

func (t *TransactionLotRelation) validate() error {
	if t.StockLotId == uuid.Nil && t.BondLotId == uuid.Nil {
		return errors.New("Both stock lot ID and bond lot ID cannot be nil")
	}
	if t.Date.After(time.Now()) {
		return errors.New("Invalid transaction date")
	}
	if t.Quantity <= 0 {
		return errors.New("Transaction quantity must be greater than 0")
	}

	return nil
}
