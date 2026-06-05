package tranlotrelation

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Relationship between a transaction and the bond/stock lots it impacted through purchases and sales
type TransactionLotRelation struct {
	Id            uuid.UUID
	TransactionId uuid.UUID
	StockLotId    uuid.UUID
	BondLotId     uuid.UUID
	Date          time.Time
	Quantity      float64
}

func New(stockLotId, bondLotId, transactionLotId uuid.UUID, date time.Time, quantity float64) (TransactionLotRelation, error) {
	t := TransactionLotRelation{
		Id:            uuid.New(),
		TransactionId: transactionLotId,
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

	return nil
}
