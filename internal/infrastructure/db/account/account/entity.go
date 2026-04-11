package accountdb

import (
	"time"

	"github.com/google/uuid"
)

type AccountDbModel struct {
	Id              uuid.UUID `sql:"id"`
	OpeningDate     time.Time `sql:"opening_date"`
	Type            string    `sql:"type"`
	Broker          string    `sql:"broker"`
	Holder          string    `sql:"holder"`
	PrimaryCurrency string    `sql:"primary_currency"`
}
