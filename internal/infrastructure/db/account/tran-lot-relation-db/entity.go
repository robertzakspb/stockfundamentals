package tranlotrelationdb

import (
	"time"

	"github.com/google/uuid"
)

type TransactionLotRelationDb struct {
	Id         uuid.UUID `sql:"id"`
	StockLotId uuid.UUID `sql:"stock_lot_id"`
	BondLotId  uuid.UUID `sql:"bond_lot_id"`
	Date       time.Time `sql:"date"`
	Quantity   float64   `sql:"quantity"`
}
